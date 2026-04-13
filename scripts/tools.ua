# this takes the input vector from a flat array and shapes it to the correct size
ParseStdio ← ↯∞_30⋕⊜□⊸≠@ &sc

# this takes the matrix and and returns a flat matrix out of it
# this also prints out to the std
FlatMatrixOut ← &p@ ≡(&pf &pf@ )♭

# this is for getting the keys for the args
Getkeys   ← ≡(⊢ ⊜□ ⊸≠@= °□)↘1&args
Getvalues ← ≡(⊣ ⊜□ ⊸≠@= °□)↘1&args
GetArgs   ← map Getkeys Getvalues

#matrix offsets 
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

#helper functions for using to get data (useful for pruning down files and doing other logic)
GetImmutableRows ← ▽¬≡⊡IsMutableIndex⟜∘
GetRowsWithRruls ← ▽ ≡>1 ≡⊡FrequencyIndex⟜∘
GetFutureRows    ← ▽≡(=4/+>(↙4datetime now) (↘1↙MinuteIndex))⟜∘
