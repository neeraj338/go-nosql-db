package command

import (
	"testing"
)

func TestValidateIfNoArgsSupplied(t *testing.T) {
	args := Args{}
	validator := CommandLineArgValidator{CommandLineArgs: args}
	err := validator.Validate()
	if err == nil {
		t.Errorf("should be an error")
	}
}

func TestValidateEither_SelectOrDelete(t *testing.T) {
	args := Args{Select: "name", Delete: true}
	validator := CommandLineArgValidator{CommandLineArgs: args}
	err := validator.Validate()
	if err == nil {
		t.Errorf("should be an error")
	}
}

func TestValidateEither_SelectOrInsert(t *testing.T) {
	args := Args{Select: "name", Data: "{}"}
	validator := CommandLineArgValidator{CommandLineArgs: args}
	err := validator.Validate()
	if err == nil {
		t.Errorf("should be an error")
	}
}

func TestValidateEither_InsertOrDelete(t *testing.T) {
	args := Args{Delete: true, Data: "{}"}
	validator := CommandLineArgValidator{CommandLineArgs: args}
	err := validator.Validate()
	if err == nil {
		t.Errorf("should be an error")
	}
}

func TestValidateEither_Select_Delete_Or_Insert(t *testing.T) {
	args := Args{Delete: true, Data: "{}", Select: "name"}
	validator := CommandLineArgValidator{CommandLineArgs: args}
	err := validator.Validate()
	if err == nil {
		t.Errorf("should be an error")
	}
}

func TestValidateAllowIfOnlyFilterPresent(t *testing.T) {
	args := Args{FilterKey: "key", FilterValue: "test"}
	validator := CommandLineArgValidator{CommandLineArgs: args}
	err := validator.Validate()
	if err != nil {
		t.Errorf("should Not be an error")
	}
}

func TestValidateAllowIfOnlyFilterValuePresent(t *testing.T) {
	args := Args{FilterValue: "test"}
	validator := CommandLineArgValidator{CommandLineArgs: args}
	err := validator.Validate()
	if err != nil {
		t.Errorf("should Not be an error")
	}
}

func TestDeleteMustHaveFilterArgs_OnlyValue_shouldAllow(t *testing.T) {
	args := Args{Delete: true, FilterValue: "test"}
	deleteCmdValidator := CommandLineArgValidator{CommandLineArgs: args}
	err := deleteCmdValidator.Validate()
	if err != nil {
		t.Errorf("should Not be an error")
	}
}

func TestDeleteMustHaveFilterArgs_OnlyKey_shouldFail(t *testing.T) {
	args := Args{Delete: true, FilterKey: "test"}
	deleteCmdValidator := CommandLineArgValidator{CommandLineArgs: args}
	err := deleteCmdValidator.Validate()
	if err == nil {
		t.Errorf("should be an error")
	}
}

func TestDeleteMustHaveFilterArgs_Both_KayAndValue(t *testing.T) {
	args := Args{Delete: true, FilterKey: "test", FilterValue: "value"}
	deleteCmdValidator := CommandLineArgValidator{CommandLineArgs: args}
	err := deleteCmdValidator.Validate()
	if err != nil {
		t.Errorf("should Not be an error")
	}
}
