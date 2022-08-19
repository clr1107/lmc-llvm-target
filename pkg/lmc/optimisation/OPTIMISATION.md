# Optimisation
Static optimisation that will read through a program and, in place, optimise its instructions. Single optimisation methods can be applied, or many chained together.

## Methods of optimisation
- (STA LDA, Thrashing.) Reducing LDA and STA instructions.
- (Remove unused boxes, Waste.) Remove unused boxes; i.e., those that are not used.
- (Box propagation, BProp.) Remove unnecessary boxes by changing which boxes instructions use.
- (Addition chaining, Chaining.) Not sure how yet lol.
- (Unrolling loop, Unroll.) I can't yet see a case where this isn't ideal if a constant is known.