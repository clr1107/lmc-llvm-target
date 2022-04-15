## Methods of optimisation
- (STA LDA.) Replace storing then loading in the same box. If no LDA or INP instructions before next STA, or at all, delete the instruction pair altogether. Otherwise, replace with just one STA.
- (Remove unused boxes.) If a box is declared yet not once used by any instructions, remove it. Apply this after every pass of optimiser.
- (Box propogation.) If a box is loaded, then stored in another, and the original was a constant, or not yet mentioned, and the dst was not yet mentioned, replace with initialised box.
- (Addition chaining.) See Stef.
- (Unrolling loop.) I can't yet see a case where this isn't ideal if a constant is known.