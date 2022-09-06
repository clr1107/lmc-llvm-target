#include "lmc.h"

void _lmc(void)
{
    int a;         // A DAT 0
    input(&a);     // INP ; STA A

    a *= 3;        // Multiplies by 3
    output(&a);    // LDA A ; OUT A
}
