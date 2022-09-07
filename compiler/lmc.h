#ifndef _LMC_BUILTIN_H
#define _LMC_BUILTIN_H

// Types available for use
typedef int number_t;
typedef int bool_t;

/**
 * The NULL pointer points to the temp mailbox, if there is one. If it does not
 * exist, this will cause it to be created.
 */
#define _TEMP ((number_t *) 0)
#define true ((bool_t) 1)
#define false ((bool_t) 0)

// Entry point.
void _lmc(void);

/**
 * `HLT` instruction.
 *
 * [Builtin]
 */
extern void _hlt(void);

/**
 * `INP` instruction.
 *
 * [Builtin]
 */
extern void _inp(void);

/**
 * `OUT` instruction.
 *
 * [Builtin]
 */
extern void _out(void);

/**
 * `STA` instruction.
 *
 * [Builtin]
 */
extern void _sta(number_t *);

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