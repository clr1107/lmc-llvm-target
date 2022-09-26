#ifndef _LMC_BUILTIN_H
#define _LMC_BUILTIN_H

// General useful macros
#define assert_int_constant(x) ((void)sizeof(struct {int dummy: 1 + !(x);}))

// Types available for use
typedef int number_t;
typedef int bool_t;

// The NULL pointer points to the temp mailbox.
#define _TEMP ((number_t *) 0)
#define true ((bool_t) 1)
#define false ((bool_t) 0)

// Compiler options and attributes
extern void __lmc_option__(const char *, number_t);
#define __lmc_option__(k, v) (assert_int_constant((v)), __lmc_option__((k), (v)))
#define __lmc_attribute__(v) __attribute__ ((annotate("___lmc__" v "___")))

#define O_NONE      0
#define O_THRASHING 1
#define O_CLEAN     2
#define O_BPROP     4
#define O_ALL       7

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