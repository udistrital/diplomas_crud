package services

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/udistrital/diplomas_crud/models"
)

// Todas las escrituras de asignaciones activas toman primero la relación y
// después la ceremonia. Inactivación y asignación se serializan entre instancias
// del CRUD usando las filas existentes, sin triggers ni cambios de esquema.
func cambiarCeremoniaTx(db orm.Ormer, cambio func(orm.TxOrmer) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := cambio(tx); err != nil {
		return err
	}
	return tx.Commit()
}

func bloquearCeremoniaTx(tx orm.TxOrmer, id int64) (bool, error) {
	var activa bool
	err := tx.Raw("SELECT activo FROM ceremonia_grado WHERE id = ? FOR UPDATE", id).QueryRow(&activa)
	return activa, err
}

func validarCeremoniaSinEstudiantesTx(tx orm.TxOrmer, id int64) error {
	var ocupada bool
	err := tx.Raw(`SELECT EXISTS (
  SELECT 1 FROM ceremonia_estudiante ce
  JOIN ceremonia_facultad_programa cfp ON cfp.id=ce.ceremonia_facultad_programa_id
  WHERE cfp.ceremonia_grado_id=? AND ce.activo IS TRUE
 )`, id).QueryRow(&ocupada)
	if err != nil {
		return err
	}
	if ocupada {
		return fmt.Errorf("%w: no se puede inactivar una ceremonia con estudiantes asignados", ErrCeremoniaConflict)
	}
	return nil
}

func validarRelacionSinEstudiantesTx(tx orm.TxOrmer, id int64) error {
	var ocupada bool
	err := tx.Raw("SELECT EXISTS (SELECT 1 FROM ceremonia_estudiante WHERE ceremonia_facultad_programa_id=? AND activo IS TRUE)", id).QueryRow(&ocupada)
	if err != nil {
		return err
	}
	if ocupada {
		return fmt.Errorf("%w: no se puede inactivar o cambiar una relación con estudiantes asignados", ErrCeremoniaConflict)
	}
	return nil
}

func bloquearAsignacionActivaTx(tx orm.TxOrmer, relacionID int64) error {
	var ceremoniaID int64
	var activa bool
	if err := tx.Raw("SELECT ceremonia_grado_id, activo FROM ceremonia_facultad_programa WHERE id=? FOR UPDATE", relacionID).QueryRow(&ceremoniaID, &activa); err != nil {
		return err
	}
	if !activa {
		return fmt.Errorf("%w: no se pueden asignar estudiantes a una relación inactiva", ErrCeremoniaConflict)
	}
	activa, err := bloquearCeremoniaTx(tx, ceremoniaID)
	if err != nil {
		return err
	}
	if !activa {
		return fmt.Errorf("%w: no se pueden asignar estudiantes a una ceremonia inactiva", ErrCeremoniaConflict)
	}
	return nil
}

func (s CeremoniaGradoService) Update(id int64, input *models.CeremoniaGrado) (*models.CeremoniaGrado, error) {
	if input == nil {
		return nil, ErrCeremoniaConflict
	}
	input.Id = id
	err := cambiarCeremoniaTx(s.ormer(), func(tx orm.TxOrmer) error {
		if _, err := bloquearCeremoniaTx(tx, id); err != nil {
			return err
		}
		if !input.Activo {
			if err := validarCeremoniaSinEstudiantesTx(tx, id); err != nil {
				return err
			}
		}
		_, err := tx.Update(input)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s CeremoniaGradoService) SoftDelete(id int64) error {
	return cambiarCeremoniaTx(s.ormer(), func(tx orm.TxOrmer) error {
		if _, err := bloquearCeremoniaTx(tx, id); err != nil {
			return err
		}
		if err := validarCeremoniaSinEstudiantesTx(tx, id); err != nil {
			return err
		}
		_, err := tx.Raw("UPDATE ceremonia_grado SET activo=false, fecha_modificacion=now() WHERE id=?", id).Exec()
		return err
	})
}

func (s CeremoniaFacultadProgramaService) Update(id int64, input *models.CeremoniaFacultadPrograma) (*models.CeremoniaFacultadPrograma, error) {
	if input == nil || input.CeremoniaGrado == nil {
		return nil, ErrCeremoniaConflict
	}
	input.Id = id
	err := cambiarCeremoniaTx(s.ormer(), func(tx orm.TxOrmer) error {
		var anterior int64
		if err := tx.Raw("SELECT ceremonia_grado_id FROM ceremonia_facultad_programa WHERE id=? FOR UPDATE", id).QueryRow(&anterior); err != nil {
			return err
		}
		if !input.Activo || anterior != input.CeremoniaGrado.Id {
			if err := validarRelacionSinEstudiantesTx(tx, id); err != nil {
				return err
			}
		}
		if input.Activo {
			activa, err := bloquearCeremoniaTx(tx, input.CeremoniaGrado.Id)
			if err != nil {
				return err
			}
			if !activa {
				return fmt.Errorf("%w: la ceremonia está inactiva", ErrCeremoniaConflict)
			}
		}
		_, err := tx.Update(input)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

func (s CeremoniaFacultadProgramaService) SoftDelete(id int64) error {
	return cambiarCeremoniaTx(s.ormer(), func(tx orm.TxOrmer) error {
		var actual int64
		if err := tx.Raw("SELECT id FROM ceremonia_facultad_programa WHERE id=? FOR UPDATE", id).QueryRow(&actual); err != nil {
			return err
		}
		if err := validarRelacionSinEstudiantesTx(tx, id); err != nil {
			return err
		}
		_, err := tx.Raw("UPDATE ceremonia_facultad_programa SET activo=false, fecha_modificacion=now() WHERE id=?", id).Exec()
		return err
	})
}

func (s CeremoniaEstudianteService) Create(input *models.CeremoniaEstudiante) (*models.CeremoniaEstudiante, error) {
	if input == nil || input.CeremoniaFacultadPrograma == nil {
		return nil, ErrCeremoniaConflict
	}
	input.Activo = true
	err := cambiarCeremoniaTx(s.ormer(), func(tx orm.TxOrmer) error {
		if err := bloquearAsignacionActivaTx(tx, input.CeremoniaFacultadPrograma.Id); err != nil {
			return err
		}
		_, err := tx.Insert(input)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(input.Id)
}

func (s CeremoniaEstudianteService) Update(id int64, input *models.CeremoniaEstudiante) (*models.CeremoniaEstudiante, error) {
	if input == nil || input.CeremoniaFacultadPrograma == nil {
		return nil, ErrCeremoniaConflict
	}
	input.Id = id
	err := cambiarCeremoniaTx(s.ormer(), func(tx orm.TxOrmer) error {
		var actual int64
		if err := tx.Raw("SELECT id FROM ceremonia_estudiante WHERE id=? FOR UPDATE", id).QueryRow(&actual); err != nil {
			return err
		}
		if input.Activo {
			if err := bloquearAsignacionActivaTx(tx, input.CeremoniaFacultadPrograma.Id); err != nil {
				return err
			}
		}
		_, err := tx.Update(input)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.GetByID(id)
}
