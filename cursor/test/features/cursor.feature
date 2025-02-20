Feature: Cursor
  Scenario: Put key-value records and then get individually
    Given there is a new tree "tree"
    When I open a new cursor "cursor" at the root of "tree"
    When I put "key09", "value09" using "cursor"
    When I put "key19", "value19" using "cursor"
    When I put "key10", "value10" using "cursor"
    When I put "key04", "value04" using "cursor"
    When I put "key21", "value21" using "cursor"
    When I put "key22", "value22" using "cursor"
    When I put "key18", "value18" using "cursor"
    When I put "key05", "value05" using "cursor"
    When I put "key23", "value23" using "cursor"
    When I put "key17", "value17" using "cursor"
    When I put "key07", "value07" using "cursor"
    When I put "key13", "value13" using "cursor"
    When I put "key24", "value24" using "cursor"
    When I put "key01", "value01" using "cursor"
    When I put "key03", "value03" using "cursor"
    When I put "key15", "value15" using "cursor"
    When I put "key08", "value08" using "cursor"
    When I put "key11", "value11" using "cursor"
    When I put "key20", "value20" using "cursor"
    When I put "key02", "value02" using "cursor"
    When I put "key12", "value12" using "cursor"
    When I put "key14", "value14" using "cursor"
    When I put "key16", "value16" using "cursor"
    When I put "key00", "value00" using "cursor"
    When I put "key06", "value06" using "cursor"
    Then I should get "key00", "value00" using "cursor"
    Then I should get "key01", "value01" using "cursor"
    Then I should get "key02", "value02" using "cursor"
    Then I should get "key03", "value03" using "cursor"
    Then I should get "key04", "value04" using "cursor"
    Then I should get "key05", "value05" using "cursor"
    Then I should get "key06", "value06" using "cursor"
    Then I should get "key07", "value07" using "cursor"
    Then I should get "key08", "value08" using "cursor"
    Then I should get "key09", "value09" using "cursor"
    Then I should get "key10", "value10" using "cursor"
    Then I should get "key11", "value11" using "cursor"
    Then I should get "key12", "value12" using "cursor"
    Then I should get "key13", "value13" using "cursor"
    Then I should get "key14", "value14" using "cursor"
    Then I should get "key15", "value15" using "cursor"
    Then I should get "key16", "value16" using "cursor"
    Then I should get "key17", "value17" using "cursor"
    Then I should get "key18", "value18" using "cursor"
    Then I should get "key19", "value19" using "cursor"
    Then I should get "key20", "value20" using "cursor"
    Then I should get "key21", "value21" using "cursor"
    Then I should get "key22", "value22" using "cursor"
    Then I should get "key23", "value23" using "cursor"
    Then I should get "key24", "value24" using "cursor"

  Scenario: Put key-value records and then get in lexicographic order
    Given there is a new tree "tree"
    When I open a new cursor "cursor" at the root of "tree"
    When I put "key13", "value13" using "cursor"
    When I put "key21", "value21" using "cursor"
    When I put "key16", "value16" using "cursor"
    When I put "key17", "value17" using "cursor"
    When I put "key04", "value04" using "cursor"
    When I put "key12", "value12" using "cursor"
    When I put "key15", "value15" using "cursor"
    When I put "key19", "value19" using "cursor"
    When I put "key01", "value01" using "cursor"
    When I put "key10", "value10" using "cursor"
    When I put "key11", "value11" using "cursor"
    When I put "key02", "value02" using "cursor"
    When I put "key22", "value22" using "cursor"
    When I put "key24", "value24" using "cursor"
    When I put "key08", "value08" using "cursor"
    When I put "key03", "value03" using "cursor"
    When I put "key09", "value09" using "cursor"
    When I put "key14", "value14" using "cursor"
    When I put "key06", "value06" using "cursor"
    When I put "key20", "value20" using "cursor"
    When I put "key23", "value23" using "cursor"
    When I put "key00", "value00" using "cursor"
    When I put "key07", "value07" using "cursor"
    When I put "key05", "value05" using "cursor"
    When I put "key18", "value18" using "cursor"
    Then I should get "key00", "value00" first using "cursor"
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

  Scenario: Put key-value records and then get in reverse lexicographic order
    Given there is a new tree "tree"
    When I open a new cursor "cursor" at the root of "tree"
    When I put "key24", "value24" using "cursor"
    When I put "key13", "value13" using "cursor"
    When I put "key01", "value01" using "cursor"
    When I put "key21", "value21" using "cursor"
    When I put "key06", "value06" using "cursor"
    When I put "key03", "value03" using "cursor"
    When I put "key22", "value22" using "cursor"
    When I put "key04", "value04" using "cursor"
    When I put "key08", "value08" using "cursor"
    When I put "key23", "value23" using "cursor"
    When I put "key02", "value02" using "cursor"
    When I put "key15", "value15" using "cursor"
    When I put "key17", "value17" using "cursor"
    When I put "key05", "value05" using "cursor"
    When I put "key07", "value07" using "cursor"
    When I put "key09", "value09" using "cursor"
    When I put "key16", "value16" using "cursor"
    When I put "key10", "value10" using "cursor"
    When I put "key12", "value12" using "cursor"
    When I put "key14", "value14" using "cursor"
    When I put "key00", "value00" using "cursor"
    When I put "key11", "value11" using "cursor"
    When I put "key19", "value19" using "cursor"
    When I put "key18", "value18" using "cursor"
    When I put "key20", "value20" using "cursor"
    Then I should get "key24", "value24" first using "cursor" in reverse
    Then I should get "key23", "value23" next using "cursor" in reverse
    Then I should get "key22", "value22" next using "cursor" in reverse
    Then I should get "key21", "value21" next using "cursor" in reverse
    Then I should get "key20", "value20" next using "cursor" in reverse
    Then I should get "key19", "value19" next using "cursor" in reverse
    Then I should get "key18", "value18" next using "cursor" in reverse
    Then I should get "key17", "value17" next using "cursor" in reverse
    Then I should get "key16", "value16" next using "cursor" in reverse
    Then I should get "key15", "value15" next using "cursor" in reverse
    Then I should get "key14", "value14" next using "cursor" in reverse
    Then I should get "key13", "value13" next using "cursor" in reverse
    Then I should get "key12", "value12" next using "cursor" in reverse
    Then I should get "key11", "value11" next using "cursor" in reverse
    Then I should get "key10", "value10" next using "cursor" in reverse
    Then I should get "key09", "value09" next using "cursor" in reverse
    Then I should get "key08", "value08" next using "cursor" in reverse
    Then I should get "key07", "value07" next using "cursor" in reverse
    Then I should get "key06", "value06" next using "cursor" in reverse
    Then I should get "key05", "value05" next using "cursor" in reverse
    Then I should get "key04", "value04" next using "cursor" in reverse
    Then I should get "key03", "value03" next using "cursor" in reverse
    Then I should get "key02", "value02" next using "cursor" in reverse
    Then I should get "key01", "value01" next using "cursor" in reverse
    Then I should get "key00", "value00" next using "cursor" in reverse
    Then getting next using "cursor" in reverse should not find

  Scenario: Get key-value records in (rev.) lex/c order, starting in the middle
    Given there is a new tree "tree"
    When I open a new cursor "cursor" at the root of "tree"
    When I put "key00", "value00" using "cursor"
    When I put "key01", "value01" using "cursor"
    When I put "key02", "value02" using "cursor"
    When I put "key03", "value03" using "cursor"
    When I put "key04", "value04" using "cursor"
    When I put "key05", "value05" using "cursor"
    When I put "key06", "value06" using "cursor"
    When I put "key07", "value07" using "cursor"
    When I put "key08", "value08" using "cursor"
    When I put "key09", "value09" using "cursor"
    When I put "key10", "value10" using "cursor"
    When I put "key11", "value11" using "cursor"
    When I put "key12", "value12" using "cursor"
    When I put "key13", "value13" using "cursor"
    When I put "key14", "value14" using "cursor"
    When I put "key15", "value15" using "cursor"
    When I put "key16", "value16" using "cursor"
    When I put "key17", "value17" using "cursor"
    When I put "key18", "value18" using "cursor"
    When I put "key19", "value19" using "cursor"
    When I put "key20", "value20" using "cursor"
    When I put "key21", "value21" using "cursor"
    When I put "key22", "value22" using "cursor"
    When I put "key23", "value23" using "cursor"
    When I put "key24", "value24" using "cursor"
    Then I should get "key12", "value12" using "cursor"
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
    Then I should get "key12", "value12" using "cursor"
    Then I should get "key11", "value11" next using "cursor" in reverse
    Then I should get "key10", "value10" next using "cursor" in reverse
    Then I should get "key09", "value09" next using "cursor" in reverse
    Then I should get "key08", "value08" next using "cursor" in reverse
    Then I should get "key07", "value07" next using "cursor" in reverse
    Then I should get "key06", "value06" next using "cursor" in reverse
    Then I should get "key05", "value05" next using "cursor" in reverse
    Then I should get "key04", "value04" next using "cursor" in reverse
    Then I should get "key03", "value03" next using "cursor" in reverse
    Then I should get "key02", "value02" next using "cursor" in reverse
    Then I should get "key01", "value01" next using "cursor" in reverse
    Then I should get "key00", "value00" next using "cursor" in reverse
    Then getting next using "cursor" in reverse should not find

  Scenario: Get key-value records in (rev.) lex/c order, changing directions
    Given there is a new tree "tree"
    When I open a new cursor "cursor" at the root of "tree"
    When I put "key02", "value02" using "cursor"
    When I put "key01", "value01" using "cursor"
    When I put "key09", "value09" using "cursor"
    When I put "key08", "value08" using "cursor"
    When I put "key05", "value05" using "cursor"
    When I put "key00", "value00" using "cursor"
    When I put "key03", "value03" using "cursor"
    When I put "key04", "value04" using "cursor"
    When I put "key06", "value06" using "cursor"
    When I put "key07", "value07" using "cursor"
    Then I should get "key08", "value08" next using "cursor"
    Then I should get "key09", "value09" next using "cursor"
    Then getting next using "cursor" should not find
    Then I should get "key09", "value09" next using "cursor" in reverse
    Then I should get "key08", "value08" next using "cursor" in reverse
    Then I should get "key07", "value07" next using "cursor" in reverse
    Then I should get "key06", "value06" next using "cursor" in reverse
    Then I should get "key05", "value05" next using "cursor" in reverse
    Then I should get "key04", "value04" next using "cursor" in reverse
    Then I should get "key03", "value03" next using "cursor" in reverse
    Then I should get "key02", "value02" next using "cursor" in reverse
    Then I should get "key01", "value01" next using "cursor" in reverse
    Then I should get "key00", "value00" next using "cursor" in reverse
    Then getting next using "cursor" in reverse should not find
    Then I should get "key00", "value00" next using "cursor"
    Then I should get "key01", "value01" next using "cursor"
    Then I should get "key02", "value02" next using "cursor"
    Then I should get "key03", "value03" next using "cursor"
    Then I should get "key04", "value04" next using "cursor"
    Then I should get "key05", "value05" next using "cursor"
    Then I should get "key06", "value06" next using "cursor"
    Then I should get "key07", "value07" next using "cursor"

  Scenario: Delete, restore, replace, and overwrite key-value records, then get
    Given there is a new tree "tree"
    When I open a new cursor "cursor" at the root of "tree"
    When I put "key05", "value05" using "cursor"
    When I put "key06", "value06" using "cursor"
    When I put "key00", "value00" using "cursor"
    When I put "key08", "value08" using "cursor"
    When I put "key07", "value07" using "cursor"
    When I put "key04", "value04" using "cursor"
    When I put "key01", "value01" using "cursor"
    When I put "key09", "value09" using "cursor"
    When I put "key02", "value02" using "cursor"
    When I put "key03", "value03" using "cursor"
    When I delete using "cursor"
    Then I should get "key04", "value04" next using "cursor"
    Then I should get "key05", "value05" next using "cursor"
    Then I should get "key06", "value06" next using "cursor"
    When I delete using "cursor"
    When I put "key06", "value06" using "cursor"
    Then I should get "key07", "value07" next using "cursor"
    Then I should get "key08", "value08" next using "cursor"
    Then I should get "key09", "value09" next using "cursor"
    When I delete using "cursor"
    When I put "key09", "VALUE09" using "cursor"
    Then getting next using "cursor" should not find
    Then I should get "key09", "VALUE09" next using "cursor" in reverse
    Then I should get "key08", "value08" next using "cursor" in reverse
    Then I should get "key07", "value07" next using "cursor" in reverse
    Then I should get "key06", "value06" next using "cursor" in reverse
    Then I should get "key05", "value05" next using "cursor" in reverse
    Then I should get "key04", "value04" next using "cursor" in reverse
    Then I should get "key02", "value02" next using "cursor" in reverse
    Then I should get "key01", "value01" next using "cursor" in reverse
    Then I should get "key00", "value00" next using "cursor" in reverse
    Then getting next using "cursor" in reverse should not find
    When I put "key00", "VALUE00" using "cursor"
    Then I should get "key00", "VALUE00" using "cursor"
    Then I should get "key01", "value01" next using "cursor"
    Then I should get "key02", "value02" next using "cursor"
    Then I should get "key04", "value04" next using "cursor"
    Then I should get "key05", "value05" next using "cursor"
    Then I should get "key06", "value06" next using "cursor"
    Then I should get "key07", "value07" next using "cursor"
    Then I should get "key08", "value08" next using "cursor"
    Then I should get "key09", "VALUE09" next using "cursor"
    Then getting next using "cursor" should not find
