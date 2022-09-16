
# Optimisation

Static optimisation that will read through a program and, in place, optimise its instructions. Single optimisation  
methods can be applied, or many chained together, using the Stacking optimiser. Beware, order may matter for the  
optimisations when stacking.

## Example in code

*TODO*

## Methods of optimisation

Explanation and files for various algorithms implemented. Some are generic compiler algorithms, and some are  
tailored to LMC.

### Clean

A simple optimisation strategy aimed at cleaning up redundant memory operations. There are two stages: 'Dead' and  
'Multi', executed in that order. This optimisation should be run after every other optimisation in the stacking format.  
This is as it cleans up redundant memory operations, commonly left by other optimisations.

#### Dead

This removes any `DAT` instructions for boxes no longer used.

E.g.,

```  
INP  
STA A  
  
A DAT 0  
B DAT 0  ; The variable 'B' is not used, this instruction can be removed.  
```  

#### Multi

This removes any `DAT` instructions for boxes already defined.

E.g.,

```  
INP  
STA A  
  
A DAT 0  
A DAT 0  ; This box is already defined above. This is redundant.  
```  

### Thrashing

This takes its name from disk thrashing; the aim of this optimisation is to reduce the use of `LDA` and `STA`  
instructions by removing redundant pairs, ineffectual instructions, etc. There are two stages: 'Multiple loading' and  
'Pairs'. They are executed in that order.

#### Multiple loading

This stage removes unused load instructions; if a mailbox is loaded yet not used (i.e., no accumulating instructions)  
before the next load instruction, the original load instruction is removed.

E.g.,

```  
INP STA A    ; Storing the input in box A  
INP  
STA B    ; Storing the second input in box B  
LDA A    ; Loading A  
LDA B    ; Loading B, but nothing happened since we loaded A (no acc  
 ;   instructions) so the original instruction, LDA A, can be removed  A DAT 0  
B DAT 0  
```  

#### Pairs

The aim of pairs is to find pairs of store and load instructions (can be similar, i.e.., store and store, or  
load and load) operating on the **same** box. If the instructions inbetween are not accumulating then the second of the  
pair can be removed.

E.g.,

```  
INP  
STA A    ; Storing the input in box A  
OUT      ; Outputting the acc's value. Does nothing to the value of the acc  
LDA A    ; Loading A again, yet nothing changed since the value was stored in A.  
 ;    Therefore, this instruction is ineffectual  A DAT 0  
```  

### Propagation

Allows the removal of unnecessary boxes by changing which boxes instructions use. I.e., find any boxes that merely serve  
a temporary purpose and remove their use, replacing them with their permanent box. The cleaning strategy can then remove  
the boxes from `DAT` instructions and thrashing can remove ineffectual loading etc. There are two stages: 'Tree', and '  
LDA STA', executed in that order.

#### Tree
A tree is computed where being the child of a parent indicates that the child derives its value from the parent. I.e.,  
for all store instructions, the source (an `INP` or `LDA` instruction with its box) is added as its parent to the global  
tree if and only if there are no accumulating instructions between the source and the `STA`.

This tree is then 'propagated'. All children are replaced by their parent's box.

E.g.,  
take this program
```  
INP  
STA A  
INP  
STA B  
LDA A  
STA D  
LDA B  
STA E  
LDA D  
SUB E  
STA C  
  
A DAT 0  
B DAT 0  
C DAT 0  
D DAT 0  
E DAT 0  
```  

The tree computed from this would be:
<div align="center"><img width="15%" src="https://i.clr.is/5HI3zCGxT.png">  </div>

(Note: `C` is not included as its source was `LDA D` but there was the accumulating `SUB E` before `STA`).

This shows that references to `E` should be replaced with `B`, and references to `D` with `A`. This would then go  
through cleaning to remove newly redundant loadings etc.

#### LDA STA
This removes instances of loading a box, then storing it in the same location without any accumulating instructions  
in between, rendering it ineffectual.

E.g.,
```  
LDA A  
OUT  
STA A    ; Nothing has happened since A was loaded.  
  
A DAT 0  
```
