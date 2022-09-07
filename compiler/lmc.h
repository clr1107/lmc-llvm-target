#ifndef _LMC_BUILTIN_H
#define _LMC_BUILTIN_H

// Types available for use
typedef int number_t;

/**
 * The NULL pointer points to the temp mailbox, if there is one. If it does not
 * exist, this will cause it to be created.
 */
#define _TEMP ((number_t *) 0)

// Entry point.
void _lmc(void);

/**
 * Halt the program. Inserts an `HLT` instruction.
 *
 * [Builtin]
 */
extern void _hlt(void);

/**
 * Output the value in a mailbox.
 *
 * [Builtin]
 */
extern void output(number_t *);

/**
 * Take an input, to be stored in a mailbox.
 *
 * [Builtin]
 */
extern void input(number_t *);

#endif