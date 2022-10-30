#include <stdio.h>

int any(const char s1[], const char s2[]) 
{
    int i = 0;
    int j;
    while(s2[i] != '\0')
    {
        j = 0;
        while (s1[j] != '\0')
        {
            /* code */
            if (s2[i] == s1[j])
                return j;
            j++;
        }
        i++;
    }
    return -1;
}

int main() {
    printf("%d\n", any("abcd", "a"));
    printf("%d\n", any("abcd", "d"));
    printf("%d\n", any("abcd", "bc"));
    printf("%d\n", any("abcd", "ef"));

}