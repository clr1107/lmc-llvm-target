#ifndef _LMC_BUILTIN_H
#define _LMC_BUILTIN_H

// General useful macros
#define assert_int_constant(x) ((void)sizeof(struct {int dummy: 1 + !(x);}))

// Types available for use
typedef int _type;
typedef _type number;
typedef _type bool;

// The NULL pointer points to the temp mailbox.
#define TEMP ((_type *) 0)
#define true ((bool) 1)
#define false ((bool) 0)

// Compiler options and attributes
extern void __lmc_option__(const char *, number);
#define __lmc_option__(k, v) (assert_int_constant((v)), __lmc_option__((k), (v)))

#define O_NONE      0
#define O_THRASHING 1
#define O_CLEAN     2
#define O_BPROP     4
#define O_ALL       7

// Set the temporary mailbox to a value
#define _mem_temp_set(v)                                        \
    _Pragma("GCC diagnostic push")                              \
    _Pragma("GCC diagnostic ignored \"-Wnull-dereference\"")    \
    *TEMP = (_type) v;                                          \
    _Pragma("GCC diagnostic pop")

// Outputs a value rather than a pointer by utilising the temporary mailbox
#define put(v)          \
    _mem_temp_set(v)    \
    output(TEMP)

// Entry point.
void _lmc(void);

/**
 * `HLT` instruction.
 *
 * [Builtin]
 */
extern void hlt(void);

/**
 * `INP` instruction.
 *
 * [Builtin]
 */
extern void inp(void);

/**
 * `OUT` instruction.
 *
 * [Builtin]
 */
extern void out(void);

/**
 * `STA` instruction.
 *
 * [Builtin]
 */
extern void sta(number *);

/**
 * Output the value in a mailbox.
 *
 * [Builtin]
 */
extern void output(number *);

/**
 * Take an input, to be stored in a mailbox.
 *
 * [Builtin]
 */
extern void input(number *);

#endif