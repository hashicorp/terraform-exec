variables {
  test = "test value"
}

run "variable_output_passthrough" {
  command = apply

  assert {
    condition     = output.test == "test value"
    error_message = "variable was not passed through to output"
  }
}
