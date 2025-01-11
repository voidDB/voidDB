Feature: Void
  Scenario: Get non-existent record
    Given there is a new Void "void" of size 1048576
    When I begin a read-only transaction "txn" in "void"
    Then getting "key" from "txn" should not find

  Scenario: Put and get record
    Given there is a new Void "void" of size 1048576
    When I begin a transaction "txn0" in "void"
    When I put "key", "value" in "txn0"
    Then I should get "key", "value" from "txn0"
    When I commit "txn0"
    When I begin a read-only transaction "txn1" in "void"
    Then I should get "key", "value" from "txn1"
    When I begin a transaction "txn2" in "void"
    When I put "key", "VALUE" in "txn2"
    Then I should get "key", "value" from "txn1"
    When I commit "txn2"
    Then I should get "key", "value" from "txn1"
    When I begin a read-only transaction "txn3" in "void"
    Then I should get "key", "VALUE" from "txn3"

  Scenario: Put record but abort transaction, then get
    Given there is a new Void "void" of size 1048576
    When I begin a transaction "txn0" in "void"
    When I put "key", "value" in "txn0"
    Then I should get "key", "value" from "txn0"
    When I abort "txn0"
    When I begin a read-only transaction "txn1" in "void"
    Then getting "key" from "txn1" should not find

  Scenario: Delete and get record
    Given there is a new Void "void" of size 1048576
    When I begin a transaction "txn0" in "void"
    When I put "key", "value" in "txn0"
    When I commit "txn0"
    When I begin a transaction "txn1" in "void"
    When I begin a read-only transaction "txn2" in "void"
    Then I should get "key", "value" from "txn1"
    When I delete "key" from "txn1"
    Then getting "key" from "txn1" should not find
    When I commit "txn1"
    Then I should get "key", "value" from "txn2"
    When I begin a read-only transaction "txn3" in "void"
    Then getting "key" from "txn3" should not find

  Scenario: Delete record but abort transaction, then get
    Given there is a new Void "void" of size 1048576
    When I begin a transaction "txn0" in "void"
    When I put "key", "value" in "txn0"
    When I commit "txn0"
    When I begin a transaction "txn1" in "void"
    When I delete "key" from "txn1"
    Then getting "key" from "txn1" should not find
    When I abort "txn1"
    When I begin a read-only transaction "txn2" in "void"
    Then I should get "key", "value" from "txn2"
