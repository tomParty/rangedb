// Code generated by go generate; DO NOT EDIT.
// This file was generated at
// 2021-01-20 15:37:16.279422 -0800 PST m=+0.003476749
package cryptotest

import (
	"strconv"

	"github.com/inklabs/rangedb/pkg/crypto"
)

func (e CustomerSignedUp) AggregateID() string   { return e.ID }
func (e CustomerSignedUp) AggregateType() string { return "customer" }
func (e CustomerSignedUp) EventType() string     { return "CustomerSignedUp" }
func (e *CustomerSignedUp) Encrypt(encryptor crypto.Encryptor) error {
	var err error
	e.Name, err = encryptor.Encrypt(e.ID, e.Name)
	if err != nil {
		e.Name = ""
		e.Email = ""
		return err
	}

	e.Email, err = encryptor.Encrypt(e.ID, e.Email)
	if err != nil {
		e.Name = ""
		e.Email = ""
		return err
	}

	return nil
}
func (e *CustomerSignedUp) Decrypt(encryptor crypto.Encryptor) error {
	var err error
	e.Name, err = encryptor.Decrypt(e.ID, e.Name)
	if err != nil {
		e.Name = ""
		e.Email = ""
		return err
	}

	e.Email, err = encryptor.Decrypt(e.ID, e.Email)
	if err != nil {
		e.Name = ""
		e.Email = ""
		return err
	}

	return nil
}

func (e CustomerAddedBirth) AggregateID() string   { return e.ID }
func (e CustomerAddedBirth) AggregateType() string { return "customer" }
func (e CustomerAddedBirth) EventType() string     { return "CustomerAddedBirth" }
func (e *CustomerAddedBirth) Encrypt(encryptor crypto.Encryptor) error {
	var err error
	stringBirthMonth := strconv.Itoa(e.BirthMonth)
	e.BirthMonthEncrypted, err = encryptor.Encrypt(e.ID, stringBirthMonth)
	if err != nil {
		e.BirthMonth = 0
		e.BirthYear = 0
		return err
	}
	e.BirthMonth = 0

	stringBirthYear := strconv.Itoa(e.BirthYear)
	e.BirthYearEncrypted, err = encryptor.Encrypt(e.ID, stringBirthYear)
	if err != nil {
		e.BirthMonth = 0
		e.BirthYear = 0
		return err
	}
	e.BirthYear = 0

	return nil
}
func (e *CustomerAddedBirth) Decrypt(encryptor crypto.Encryptor) error {
	var err error
	decryptedBirthMonth, err := encryptor.Decrypt(e.ID, e.BirthMonthEncrypted)
	if err != nil {
		e.BirthMonth = 0
		e.BirthYear = 0
		return err
	}
	e.BirthMonth, err = strconv.Atoi(decryptedBirthMonth)
	if err != nil {
		e.BirthMonth = 0
		e.BirthYear = 0
		return err
	}
	e.BirthMonthEncrypted = ""

	decryptedBirthYear, err := encryptor.Decrypt(e.ID, e.BirthYearEncrypted)
	if err != nil {
		e.BirthMonth = 0
		e.BirthYear = 0
		return err
	}
	e.BirthYear, err = strconv.Atoi(decryptedBirthYear)
	if err != nil {
		e.BirthMonth = 0
		e.BirthYear = 0
		return err
	}
	e.BirthYearEncrypted = ""

	return nil
}
