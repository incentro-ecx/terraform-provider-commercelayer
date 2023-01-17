
resource "commercelayer_manual_tax_calculator" "incentro_manual_tax_calculator" {
  attributes {
    name = "Incentro Manual Tax Calculator"
  }
  relationships {
    tax_rule_id = commercelayer_tax_rule.incentro_tax_rule.id
  }
}

resource "commercelayer_tax_rule" "incentro_tax_rule" {
  attributes {
    name = "Incentro Tax Rule"
    metadata = {
      foo : "bar"
      testName : "{{.testName}}"
    }
  }
  relationships {
    manual_tax_calculator_id = commercelayer_manual_tax_calculator.incentro_manual_tax_calculator.id
  }
}
