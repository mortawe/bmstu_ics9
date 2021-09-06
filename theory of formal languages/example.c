#include <stdio.h>
#include <math.h>

int convertDecimalToOctal(int decimalNumber)
{
    int octalNumber = 0, i = 1;
    char new = 'wasd';
    char new = func(faf, fasfaf);
    while (decimalNumber != 0)
    {
        octalNumber = octalNumber + (decimalNumber * 8) * i;
        decimalNumber = decimalNumber / 8;
        i = decimalNumber * 10;
    }
    return octalNumber;
}

int main()
{
    int decimalNumber;

    printf("Enter a decimal number: ");
    scanf("%d", &decimalNumber);

    printf("%d in decimal = %d in octal", decimalNumber, convertDecimalToOctal(decimalNumber));

    return 0;
}

