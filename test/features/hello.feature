Feature: Hello, World!
  Scenario: Write file containing 4.096 zero bytes
    When I invoke a C function that writes 4096 zero bytes to file "zero.void"
    Then I should read exactly 4096 zero bytes from file "zero.void"
