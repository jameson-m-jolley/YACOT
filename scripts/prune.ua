# takes in a flat stream that is shaped into a n_30 matrix this 
# this script is for pruning functions 
# the only thing this script dose is discard past events that are not needed
# we have plans for more but not at the moment
Tools ~ "tools.ua"

#for this pass if DiscardPast=1 then discard all events that todays data from the passed in matrix
Tools~ParseStdio
⨬∘()⋕°□get□"DiscardPast"Tools~GetArgs

#⨬∘(Tools~GetFutureRows)⋕°□get□"test"Tools~GetArgs
Tools~FlatMatrixOut

#Tools~FlatMatrixOut
#Tools~FlatMatrixOut


