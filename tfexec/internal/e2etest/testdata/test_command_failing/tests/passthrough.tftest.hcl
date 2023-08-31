variables {
  test = "test value"
}

run "variable_output_passthrough" {
  command = apply

  assert {
    condition     = output.test == "not test value" # intentionally incorrect
    error_message = "variable was not passed through to output"
  }
}
