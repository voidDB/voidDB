Feature: Tree
  Scenario: Put key-value record, and then get
    Given there is a new tree "tree"
    When I put "key", "value" into "tree"
    Then I should get "key", "value" from "tree"
    When I open a new cursor "cursor" at the root of "tree"
    Then I should get "key", "value" next using "cursor"
    Then getting next using "cursor" should not find

  Scenario: Get key-value record not yet put
    Given there is a new tree "tree"
    Then getting "key" from "tree" should not find
    When I open a new cursor "cursor" at the root of "tree"
    Then getting next using "cursor" should not find
    Then getting "key" using "cursor" should not find

  Scenario: Put key-value records into tree, causing overflow (2 levels)
    Given there is a new tree "tree"
    When I put "key0", "value0" into "tree"
    When I put "key1", "value1" into "tree"
    When I put "key2", "value2" into "tree"
    When I put "key3", "value3" into "tree"
    When I put "key4", "value4" into "tree"
    When I put "key5", "value5" into "tree"
    When I put "key6", "value6" into "tree"
    Then I should get "key0", "value0" from "tree"
    Then I should get "key1", "value1" from "tree"
    Then I should get "key2", "value2" from "tree"
    Then I should get "key3", "value3" from "tree"
    Then I should get "key4", "value4" from "tree"
    Then I should get "key5", "value5" from "tree"
    Then I should get "key6", "value6" from "tree"
    When I open a new cursor "cursor" at the root of "tree"
    Then I should get "key0", "value0" next using "cursor"
    Then I should get "key1", "value1" next using "cursor"
    Then I should get "key2", "value2" next using "cursor"
    Then I should get "key3", "value3" next using "cursor"
    Then I should get "key4", "value4" next using "cursor"
    Then I should get "key5", "value5" next using "cursor"
    Then I should get "key6", "value6" next using "cursor"
    Then getting next using "cursor" should not find
    Then I should get "key6", "value6" next using "cursor" in reverse
    Then I should get "key5", "value5" next using "cursor" in reverse
    Then I should get "key4", "value4" next using "cursor" in reverse
    Then I should get "key3", "value3" next using "cursor" in reverse
    Then I should get "key2", "value2" next using "cursor" in reverse
    Then I should get "key1", "value1" next using "cursor" in reverse
    Then I should get "key0", "value0" next using "cursor" in reverse
    Then getting next using "cursor" in reverse should not find

  Scenario: Put key-value records into tree, causing overflow (3 levels)
    Given there is a new tree "tree"
    When I put "key00", "value00" into "tree"
    When I put "key01", "value01" into "tree"
    When I put "key02", "value02" into "tree"
    When I put "key03", "value03" into "tree"
    When I put "key04", "value04" into "tree"
    When I put "key05", "value05" into "tree"
    When I put "key06", "value06" into "tree"
    When I put "key07", "value07" into "tree"
    When I put "key08", "value08" into "tree"
    When I put "key09", "value09" into "tree"
    When I put "key10", "value10" into "tree"
    When I put "key11", "value11" into "tree"
    When I put "key12", "value12" into "tree"
    When I put "key13", "value13" into "tree"
    When I put "key14", "value14" into "tree"
    When I put "key15", "value15" into "tree"
    When I put "key16", "value16" into "tree"
    When I put "key17", "value17" into "tree"
    When I put "key18", "value18" into "tree"
    When I put "key19", "value19" into "tree"
    When I put "key20", "value20" into "tree"
    When I put "key21", "value21" into "tree"
    When I put "key22", "value22" into "tree"
    When I put "key23", "value23" into "tree"
    When I put "key24", "value24" into "tree"
    Then I should get "key00", "value00" from "tree"
    Then I should get "key01", "value01" from "tree"
    Then I should get "key02", "value02" from "tree"
    Then I should get "key03", "value03" from "tree"
    Then I should get "key04", "value04" from "tree"
    Then I should get "key05", "value05" from "tree"
    Then I should get "key06", "value06" from "tree"
    Then I should get "key07", "value07" from "tree"
    Then I should get "key08", "value08" from "tree"
    Then I should get "key09", "value09" from "tree"
    Then I should get "key10", "value10" from "tree"
    Then I should get "key11", "value11" from "tree"
    Then I should get "key12", "value12" from "tree"
    Then I should get "key13", "value13" from "tree"
    Then I should get "key14", "value14" from "tree"
    Then I should get "key15", "value15" from "tree"
    Then I should get "key16", "value16" from "tree"
    Then I should get "key17", "value17" from "tree"
    Then I should get "key18", "value18" from "tree"
    Then I should get "key19", "value19" from "tree"
    Then I should get "key20", "value20" from "tree"
    Then I should get "key21", "value21" from "tree"
    Then I should get "key22", "value22" from "tree"
    Then I should get "key23", "value23" from "tree"
    Then I should get "key24", "value24" from "tree"
    When I open a new cursor "cursor" at the root of "tree"
    Then I should get "key00", "value00" next using "cursor"
    Then I should get "key01", "value01" next using "cursor"
    Then I should get "key02", "value02" next using "cursor"
    Then I should get "key03", "value03" next using "cursor"
    Then I should get "key04", "value04" next using "cursor"
    Then I should get "key05", "value05" next using "cursor"
    Then I should get "key06", "value06" next using "cursor"
    Then I should get "key07", "value07" next using "cursor"
    Then I should get "key08", "value08" next using "cursor"
    Then I should get "key09", "value09" next using "cursor"
    Then I should get "key10", "value10" next using "cursor"
    Then I should get "key11", "value11" next using "cursor"
    Then I should get "key12", "value12" next using "cursor"
    Then I should get "key13", "value13" next using "cursor"
    Then I should get "key14", "value14" next using "cursor"
    Then I should get "key15", "value15" next using "cursor"
    Then I should get "key16", "value16" next using "cursor"
    Then I should get "key17", "value17" next using "cursor"
    Then I should get "key18", "value18" next using "cursor"
    Then I should get "key19", "value19" next using "cursor"
    Then I should get "key20", "value20" next using "cursor"
    Then I should get "key21", "value21" next using "cursor"
    Then I should get "key22", "value22" next using "cursor"
    Then I should get "key23", "value23" next using "cursor"
    Then I should get "key24", "value24" next using "cursor"
    Then getting next using "cursor" should not find

  Scenario: Delete key-value record, and then get
    Given there is a new tree "tree"
    When I put "key", "value" into "tree"
    When I delete "key" from "tree"
    Then getting "key" from "tree" should not find
    When I open a new cursor "cursor" at the root of "tree"
    Then getting next using "cursor" should not find
    Then getting "key" using "cursor" should not find

  Scenario: Delete key-value record, put it back, and then get
    Given there is a new tree "tree"
    When I put "key", "value" into "tree"
    When I delete "key" from "tree"
    When I put "key", "value" into "tree"
    Then I should get "key", "value" from "tree"
    When I open a new cursor "cursor" at the root of "tree"
    Then I should get "key", "value" next using "cursor"
    Then getting next using "cursor" should not find

  Scenario: Delete key-value record, causing underflow (take right)
    Given there is a new tree "tree"
    When I put "key0", "value0" into "tree"
    When I put "key1", "value1" into "tree"
    When I put "key2", "value2" into "tree"
    When I put "key3", "value3" into "tree"
    When I put "key4", "value4" into "tree"
    When I put "key5", "value5" into "tree"
    When I put "key6", "value6" into "tree"
    When I delete "key0" from "tree"
    Then getting "key0" from "tree" should not find
    Then I should get "key1", "value1" from "tree"
    Then I should get "key2", "value2" from "tree"
    Then I should get "key3", "value3" from "tree"
    Then I should get "key4", "value4" from "tree"
    Then I should get "key5", "value5" from "tree"
    Then I should get "key6", "value6" from "tree"
    When I open a new cursor "cursor" at the root of "tree"
    Then I should get "key1", "value1" next using "cursor"
    Then I should get "key2", "value2" next using "cursor"
    Then I should get "key3", "value3" next using "cursor"
    Then I should get "key4", "value4" next using "cursor"
    Then I should get "key5", "value5" next using "cursor"
    Then I should get "key6", "value6" next using "cursor"
    Then getting next using "cursor" should not find

  Scenario: Delete key-value record, causing underflow (take left)
    Given there is a new tree "tree"
    When I put "key00", "value00" into "tree"
    When I put "key10", "value10" into "tree"
    When I put "key20", "value20" into "tree"
    When I put "key30", "value30" into "tree"
    When I put "key40", "value40" into "tree"
    When I put "key50", "value50" into "tree"
    When I put "key60", "value60" into "tree"
    When I put "key15", "value15" into "tree"
    When I delete "key60" from "tree"
    When I delete "key50" from "tree"
    Then I should get "key00", "value00" from "tree"
    Then I should get "key10", "value10" from "tree"
    Then I should get "key15", "value15" from "tree"
    Then I should get "key20", "value20" from "tree"
    Then I should get "key30", "value30" from "tree"
    Then I should get "key40", "value40" from "tree"
    Then getting "key50" from "tree" should not find
    Then getting "key60" from "tree" should not find
    When I open a new cursor "cursor" at the root of "tree"
    Then I should get "key00", "value00" next using "cursor"
    Then I should get "key10", "value10" next using "cursor"
    Then I should get "key15", "value15" next using "cursor"
    Then I should get "key20", "value20" next using "cursor"
    Then I should get "key30", "value30" next using "cursor"
    Then I should get "key40", "value40" next using "cursor"
    Then getting next using "cursor" should not find

  Scenario: Delete key-value record, causing underflow (merge right)
    Given there is a new tree "tree"
    When I put "key0", "value0" into "tree"
    When I put "key1", "value1" into "tree"
    When I put "key2", "value2" into "tree"
    When I put "key3", "value3" into "tree"
    When I put "key4", "value4" into "tree"
    When I put "key5", "value5" into "tree"
    When I put "key6", "value6" into "tree"
    When I delete "key6" from "tree"
    When I delete "key0" from "tree"
    Then getting "key0" from "tree" should not find
    Then I should get "key1", "value1" from "tree"
    Then I should get "key2", "value2" from "tree"
    Then I should get "key3", "value3" from "tree"
    Then I should get "key4", "value4" from "tree"
    Then I should get "key5", "value5" from "tree"
    Then getting "key6" from "tree" should not find
    When I open a new cursor "cursor" at the root of "tree"
    Then I should get "key1", "value1" next using "cursor"
    Then I should get "key2", "value2" next using "cursor"
    Then I should get "key3", "value3" next using "cursor"
    Then I should get "key4", "value4" next using "cursor"
    Then I should get "key5", "value5" next using "cursor"
    Then getting next using "cursor" should not find

  Scenario: Delete key-value record, causing underflow (merge left)
    Given there is a new tree "tree"
    When I put "key0", "value0" into "tree"
    When I put "key1", "value1" into "tree"
    When I put "key2", "value2" into "tree"
    When I put "key3", "value3" into "tree"
    When I put "key4", "value4" into "tree"
    When I put "key5", "value5" into "tree"
    When I put "key6", "value6" into "tree"
    When I put "key7", "value7" into "tree"
    When I put "key8", "value8" into "tree"
    When I put "key9", "value9" into "tree"
    When I delete "key9" from "tree"
    When I delete "key4" from "tree"
    Then I should get "key0", "value0" from "tree"
    Then I should get "key1", "value1" from "tree"
    Then I should get "key2", "value2" from "tree"
    Then I should get "key3", "value3" from "tree"
    Then getting "key4" from "tree" should not find
    Then I should get "key5", "value5" from "tree"
    Then I should get "key6", "value6" from "tree"
    Then I should get "key7", "value7" from "tree"
    Then I should get "key8", "value8" from "tree"
    Then getting "key9" from "tree" should not find
    When I open a new cursor "cursor" at the root of "tree"
    Then I should get "key0", "value0" next using "cursor"
    Then I should get "key1", "value1" next using "cursor"
    Then I should get "key2", "value2" next using "cursor"
    Then I should get "key3", "value3" next using "cursor"
    Then I should get "key5", "value5" next using "cursor"
    Then I should get "key6", "value6" next using "cursor"
    Then I should get "key7", "value7" next using "cursor"
    Then I should get "key8", "value8" next using "cursor"
    Then getting next using "cursor" should not find

  Scenario: Get key-value record and subsequent records using cursor
    Given there is a new tree "tree"
    When I put "key0", "value0" into "tree"
    When I put "key1", "value1" into "tree"
    When I put "key2", "value2" into "tree"
    When I put "key3", "value3" into "tree"
    When I open a new cursor "cursor" at the root of "tree"
    Then I should get "key2", "value2" using "cursor"
    Then I should get "key3", "value3" next using "cursor"
    Then getting next using "cursor" should not find
    Then I should get "key0", "value0" using "cursor"
    Then I should get "key1", "value1" next using "cursor"

  Scenario: Delete key-value record using cursor
    Given there is a new tree "tree"
    When I put "key0", "value0" into "tree"
    When I put "key1", "value1" into "tree"
    When I put "key2", "value2" into "tree"
    When I open a new cursor "cursor" at the root of "tree"
    Then I should get "key0", "value0" next using "cursor"
    When I delete using "cursor"
    Then I should get "key1", "value1" next using "cursor"
    Then I should get "key2", "value2" next using "cursor"
    Then getting next using "cursor" should not find
    Then getting "key0" using "cursor" should not find

  Scenario: Put key-value record using cursor
    Given there is a new tree "tree"
    When I put "key0", "value0" into "tree"
    When I put "key2", "value2" into "tree"
    When I open a new cursor "cursor" at the root of "tree"
    When I put "key1", "value1" using "cursor"
    Then I should get "key2", "value2" next using "cursor"
    Then getting next using "cursor" should not find
    Then I should get "key0", "value0" using "cursor"
    Then I should get "key1", "value1" next using "cursor"

  Scenario: Get first and last key-value records using cursor
    Given there is a new tree "tree"
    When I put "key0", "value0" into "tree"
    When I put "key1", "value1" into "tree"
    When I put "key2", "value2" into "tree"
    When I put "key3", "value3" into "tree"
    When I open a new cursor "cursor" at the root of "tree"
    Then I should get "key3", "value3" first using "cursor" in reverse
    Then I should get "key2", "value2" next using "cursor" in reverse
    Then I should get "key0", "value0" first using "cursor"
    Then I should get "key1", "value1" next using "cursor"
