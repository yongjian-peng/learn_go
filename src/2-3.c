#include <stdio.h>
#include <stdbool.h>

/* convert a hex string to digit                */
int htoi(const char s[])
{
    int i = 0;
    int hex = 0;
    bool skipHeader = false; /* judge the header is 0x or 0x         */
    while (s[i] != '\0')
    {
        printf("%c\n", s[i]);
        if (!skipHeader && s[0] == '0' && (s[1] == 'x' || s[1] == 'X'))
        {
            skipHeader = true;
            if (s[2] == '\0') /* no digit after 0x or 0x              */
            {
                printf("lack of digits after 0x or 0X\n");
                hex = -1;
                break;
            }
            else
            {
                i = 2; /* skip header:0x or 0X                */
                continue;
            }
        }
        else
        {
            if (s[i] >= '0' && s[i] <= '9')
            {
                hex = hex * 16 + (s[i] - '0');
            }
            else if (s[i] >= 'a' && s[i] <= 'f')
            {
                hex = hex * 16 + (s[i] - 'a' + 10);
            }
            else if (s[i] >= 'A' && s[i] <= 'F')
            {
                hex = hex * 16 + (s[i] - 'A' + 10);
            }
            else
            {
                printf("error: illegal input: %c\n", s[i]);
                return -1;
            }
            i++;
        }
    }
    if (s[0] == '\0')
    {
        printf("no input\n");
        hex = -1;
    }
    return hex;
}

int main()
{
    printf("0x%x\n", htoi("1234"));
    printf("0x%x\n", htoi("0x1234"));
    printf("0x%x\n", htoi("0"));
    printf("0x%x\n", htoi("0x0"));
    printf("0x%x\n", htoi("0x1234ab"));
    printf("0x%x\n", htoi("0x1234AF"));
    printf("0x%x\n", htoi("0x123456789"));
    printf("%d\n", htoi("0x"));
    printf("%d\n", htoi("0s"));
    printf("%d\n", htoi("0x120xab"));
    printf("%d\n", htoi(""));
}