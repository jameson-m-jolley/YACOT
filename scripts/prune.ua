# Prune pass - pass through unchanged
# Input: CSV from xargs
# Output: CSV unchanged
Tools ~ "tools.ua"

#for this pass if DiscardPast=1 then discard all events that todays data from the passed in matrix
Tools~ParseStdio
#id,year,month,day,hour,minute,duration_minutes,priority,category_id,deadline_year,deadline_month,deadline_day,deadline_hour,deadline_minute,frequency,interval,until_year,until_month,until_day,until_hour,until_minute,by_day_mask,by_month_mask,by_monthday_mask,by_hour_mask,by_minute_mask,by_second_mask,count,is_mutable,wkst
⨬∘(
      ⊂ ⊙Tools~GetRowsWithRruls  Tools~GetRowsWithRruls ⟜∘
)⋕°□get□"DiscardPast"Tools~GetArgs

⨬∘(Tools~GetFutureRows)⋕°□get□"test"Tools~GetArgs

#Tools~FlatMatrixOut
#Tools~FlatMatrixOut


