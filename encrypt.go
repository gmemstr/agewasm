package main

import (
	"bytes"
	"io"
	"strings"
	"syscall/js"

	"filippo.io/age"
	"filippo.io/age/armor"
)

func Encrypt(this js.Value, args []js.Value) interface{} {
	output := make(map[string]interface{})
	if len(args) != 2 {
		output["error"] = "invalid arguments. expected: recipients, input"
		return output
	}
	var recipients = args[0].String()
	var input = args[1].String()
	buff := bytes.NewBuffer(nil)
	ids, err := age.ParseRecipients(strings.NewReader(recipients))
	if err != nil {
		output["error"] = err.Error()
		return output
	}
	err = encrypt(ids, strings.NewReader(input), buff, true)
	if err != nil {
		output["error"] = err.Error()
		return output
	}
	output["output"] = buff.String()
	return output
}

// encrypt internal helper
func encrypt(recipients []age.Recipient, in io.Reader, out io.Writer, withArmor bool) error {
	var a io.WriteCloser
	if withArmor {
		a = armor.NewWriter(out)
		out = a
	}
	w, err := age.Encrypt(out, recipients...)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, in); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	if a != nil {
		if err := a.Close(); err != nil {
			return err
		}
	}
	return nil
}
