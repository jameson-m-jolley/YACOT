
#matrix offsets 
#--------------------------------------------------
IdIndex              ← 0
YearIndex            ← 1
MonthIndex           ← 2
DayIndex             ← 3
HourIndex            ← 4
MinuteIndex          ← 5
DurationMinutesIndex ← 6
PriorityIndex        ← 7
CategoryIdIndex      ← 8
DeadlineYearIndex    ← 9
DeadlineMonthIndex   ← 10
DeadlineDayIndex     ← 11
DeadlineHourIndex    ← 12
DeadlineMinuteIndex  ← 13
FrequencyIndex       ← 14
IntervalIndex        ← 15
UntilYearIndex       ← 16
UntilMonthIndex      ← 17
UntilDayIndex        ← 18
UntilHourIndex       ← 19
UntilMinuteIndex     ← 20
ByDayMaskIndex       ← 21
ByMonthMaskIndex     ← 22
ByMonthdayMaskIndex  ← 23
ByHourMaskIndex      ← 24
ByMinuteMaskIndex    ← 25
BySecondMaskIndex    ← 26
CountIndex           ← 27
IsMutableIndex       ← 28
WkstIndex            ← 29
#--------------------------------------------------


#helper functions for using to get data (useful for pruning down files and doing other logic)
#--------------------------------------------------
# get all rows that are immutable 
GetImmutableRows ← ▽¬≡⊡IsMutableIndex⟜∘

# get all rows with a Rrule
GetRowsWithRrule ← ▽ ≡>1 ≡⊡FrequencyIndex⟜∘

# this gits rows with rules that repeat
GetFutureRows    ← ▽≡(=4/+<(↙4datetime now) (↘1↙MinuteIndex))⟜∘

# this gits rows with rules that repeat that have not ended
GetUnexpiredRruleRows = ▽≡(+ =4/+>(↙4datetime now)  ⊙ (=4/+<0) ⟜∘(↘UntilYearIndex↙UntilMinuteIndex))⟜∘

#get a row with an id
GetRowWithID = ∘


# this takes the input vector from a flat array and shapes it to the correct size
ParseStdio ← ↯∞_30⋕⊜□⊸≠@ &sc

# this takes the matrix and and returns a flat matrix out of it
# this also prints out to the std
FlatMatrixOut ← &p@ ≡(&pf &pf@ )♭

# this is for getting the keys for the args
Getkeys   ← ≡(⊢ ⊜□ ⊸≠@= °□)↘1&args
Getvalues ← ≡(⊣ ⊜□ ⊸≠@= °□)↘1&args
GetArgs   ← map Getkeys Getvalues
#--------------------------------------------------


#cost/objective functions 
#--------------------------------------------------


MinimumTime = ⊣⍆≡(    ↘+1IdIndex ↙+1MinuteIndex) 
MaximumTime = ⊢⍆≡(    ↘+1IdIndex ↙+1MinuteIndex)

#converts a datetime to minutes
TimeInMinute = /+ ×[525600 43800 1440 60 1]

# takes a cal matrix as an arg and returns a scalar oc the total time in min in the config
TotalMinutesInConfig =  TimeInMinute ◿ [ ∞ 12 31 23 59]  - MaximumTime ⟜MinimumTime 

TotalOccupiedTimeInMinutes =  ≡ (⊏ ↘FrequencyIndex⇡+1ByMinuteMaskIndex)

AvgFreeTimePerHour = ∘


# this function takes a list of strings options and the matrix 
# and returns a map of scalar that represents 
# the avg free time per unit in the map 
# the map vals that will return a computed val are as fallows:
# hour,
# day, 
# weak, 
# month,


#--------------------------------------------------
#tests. 





