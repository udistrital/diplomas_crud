package services

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/udistrital/diplomas_crud/models"
)

var ErrFirmaConflict = errors.New("las firmas del diploma cambiaron; actualice e intente nuevamente")

type RevertirFirmaDocumentoInput struct {
	DocumentoIdentidad int64   `json:"documento_identidad"`
	RolFirmanteId      int64   `json:"rol_firmante_id"`
	FirmasEsperadas    []int64 `json:"firmas_esperadas"`
	EstadoNuevoId      int64   `json:"estado_nuevo_id"`
}

type RevertirFirmaDocumentoResult struct {
	Anuladas int64 `json:"anuladas"`
}

// Every signature mutation takes the document lock before reading/writing signatures.
// The MID validates role order; this snapshot closes the gap between that validation
// and persistence, including concurrent requests handled by different MID instances.
func bloquearFirmasDocumentoTx(tx orm.TxOrmer, documentoID int64, esperadas []int64, permitirFinalizado bool) error {
	var uuid string
	if err := tx.Raw("SELECT COALESCE(uuid_documento::text, '') FROM documento_digital WHERE id = ? FOR UPDATE", documentoID).QueryRow(&uuid); err != nil {
		return err
	}
	if uuid != "" && !permitirFinalizado {
		return fmt.Errorf("%w: el diploma ya fue finalizado", ErrFirmaConflict)
	}
	if esperadas == nil {
		return nil
	}
	var actuales []models.FirmaDocumento
	if _, err := tx.Raw("SELECT id FROM firma_documento WHERE documento_digital_id = ? AND activo IS TRUE ORDER BY id", documentoID).QueryRows(&actuales); err != nil {
		return err
	}
	if len(actuales) != len(esperadas) {
		return ErrFirmaConflict
	}
	ids := make(map[int64]bool, len(esperadas))
	for _, id := range esperadas {
		if id <= 0 || ids[id] {
			return ErrFirmaConflict
		}
		ids[id] = true
	}
	for _, firma := range actuales {
		if !ids[firma.Id] {
			return ErrFirmaConflict
		}
	}
	return nil
}

func (s DocumentoDigitalWorkflowService) RevertirFirma(documentoID int64, input *RevertirFirmaDocumentoInput) (*RevertirFirmaDocumentoResult, error) {
	if input == nil || documentoID <= 0 || input.DocumentoIdentidad <= 0 || input.RolFirmanteId <= 0 || input.EstadoNuevoId <= 0 || len(input.FirmasEsperadas) == 0 {
		return nil, ErrFirmaConflict
	}
	tx, err := orm.NewOrm().Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if err := bloquearFirmasDocumentoTx(tx, documentoID, input.FirmasEsperadas, false); err != nil {
		return nil, err
	}
	// All signatures of this role must belong to the actor authorized by the MID.
	var firmas []models.FirmaDocumento
	if _, err := tx.Raw("SELECT id, documento_identidad FROM firma_documento WHERE documento_digital_id = ? AND rol_firmante_id = ? AND activo IS TRUE ORDER BY id", documentoID, input.RolFirmanteId).QueryRows(&firmas); err != nil {
		return nil, err
	}
	if len(firmas) == 0 {
		return nil, ErrFirmaConflict
	}
	for _, firma := range firmas {
		if firma.DocumentoIdentidad != input.DocumentoIdentidad {
			return nil, ErrFirmaConflict
		}
	}
	for _, firma := range firmas {
		if _, err := tx.Raw("UPDATE firma_documento SET activo = false, fecha_modificacion = now() WHERE id = ?", firma.Id).Exec(); err != nil {
			return nil, err
		}
		id := firma.Id
		if _, err := registrarEstadoDocumentoTx(tx, documentoID, &CambiarEstadoDocumentoInput{
			EstadoNuevoId: input.EstadoNuevoId, DocumentoIdentidad: &input.DocumentoIdentidad, RolActorId: &input.RolFirmanteId, FirmaDocumentoId: &id, Observacion: "Firma revertida por su titular desde el flujo de diplomas",
		}); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &RevertirFirmaDocumentoResult{Anuladas: int64(len(firmas))}, nil
}
