package services

import "testing"

func TestSiguienteUbicacionDiploma(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                                 string
		libroActual, ultimoFolio, ultimoActa int
		foliosPorLibro, actasPorFolio        int
		wantLibro, wantFolio, wantActa       int
	}{
		{
			name:           "inicia primer libro folio y acta",
			libroActual:    1,
			ultimoFolio:    0,
			ultimoActa:     0,
			foliosPorLibro: 500,
			actasPorFolio:  6,
			wantLibro:      1,
			wantFolio:      1,
			wantActa:       1,
		},
		{
			name:           "mantiene folio hasta completar seis actas",
			libroActual:    1,
			ultimoFolio:    1,
			ultimoActa:     5,
			foliosPorLibro: 500,
			actasPorFolio:  6,
			wantLibro:      1,
			wantFolio:      1,
			wantActa:       6,
		},
		{
			name:           "acta seis del folio uno pasa a acta uno del folio dos",
			libroActual:    1,
			ultimoFolio:    1,
			ultimoActa:     6,
			foliosPorLibro: 500,
			actasPorFolio:  6,
			wantLibro:      1,
			wantFolio:      2,
			wantActa:       1,
		},
		{
			name:           "ultimo folio completo inicia nuevo libro",
			libroActual:    1,
			ultimoFolio:    500,
			ultimoActa:     6,
			foliosPorLibro: 500,
			actasPorFolio:  6,
			wantLibro:      2,
			wantFolio:      1,
			wantActa:       1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotLibro, gotFolio, gotActa := siguienteUbicacionDiploma(
				tt.libroActual,
				tt.ultimoFolio,
				tt.ultimoActa,
				tt.foliosPorLibro,
				tt.actasPorFolio,
			)

			if gotLibro != tt.wantLibro || gotFolio != tt.wantFolio || gotActa != tt.wantActa {
				t.Fatalf(
					"expected libro=%d folio=%d acta=%d, got libro=%d folio=%d acta=%d",
					tt.wantLibro,
					tt.wantFolio,
					tt.wantActa,
					gotLibro,
					gotFolio,
					gotActa,
				)
			}
		})
	}
}
