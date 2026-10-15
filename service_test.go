package main

import "testing"

/*
Testes INDIVIDUAIS o debaixo foi feito todos esses em 1 unica função

func TestValidateService(t *testing.T) {
	service := Service{
		Name:  "Desenvolvimento de sites",
		Price: 1500,
	}

	err := validateService(service)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateServicePriceZero(t *testing.T) {
	testService := Service{Name: "Desenvolvimento de sites", Price: 0}

	errorMessagePrice := "Price must be greater than zero"

	err := validateService(testService)

	if err == nil {
		//t.Fatalf -> Caso caia aqui, ele finaliza o codigo e nao continua a execução para t.Errorf.
		t.Fatalf("Expecting an invalid price error.")
	}

	if err.Error() != errorMessagePrice {
		//%q -> é útil para imprimir strings entre aspas, deixando diferenças como espaços mais fáceis de perceber.
		t.Errorf("expected %q, got %q", errorMessagePrice, err.Error())
	}

}

func TestValidateServiceEmptyText(t *testing.T) {
	testService := Service{Name: "", Price: 200}

	errorMessageName := "Name must not be empty"

	err := validateService(testService)

	if err == nil {
		//t.Fatalf -> Caso caia aqui, ele finaliza o codigo e nao continua a execução para t.Errorf.
		t.Fatalf("Expecting an invalid empty name error.")
	}

	if err.Error() != errorMessageName {
		//%q -> é útil para imprimir strings entre aspas, deixando diferenças como espaços mais fáceis de perceber.
		t.Errorf("expected %q, got %q", errorMessageName, err.Error())
	}

}
*/

func TestValidateService(t *testing.T) {
	tests := []struct {
		name      string
		service   Service
		wantError string
	}{
		{
			name:    "valid service",
			service: Service{Name: "Sites", Price: 1500},
		},
		{
			name:      "zero price",
			service:   Service{Name: "Sites", Price: 0},
			wantError: "Price must be greater than zero",
		},
		{
			name:      "empty name",
			service:   Service{Name: "", Price: 200},
			wantError: "Name must not be empty",
		},
		{
			name:      "empty name and with spaces",
			service:   Service{Name: "           ", Price: 200},
			wantError: "Name must not be empty",
		},
	}

	for _, test := range tests {
		// t.Run cria um subteste isolado para cada cenário
		t.Run(test.name, func(t *testing.T) {

			err := validateService(test.service)

			if test.name != "valid service" {
				if err == nil {
					//t.Fatalf -> Caso caia aqui, ele finaliza o codigo e nao continua a execução para t.Errorf.
					t.Fatalf("Expecting an invalid error. %s", test.wantError)
				}

				if err.Error() != test.wantError {
					//%q -> é útil para imprimir strings entre aspas, deixando diferenças como espaços mais fáceis de perceber.
					t.Errorf("expected %q, got %q", test.wantError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}

}
