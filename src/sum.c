#include <stdio.h>

void fun(char*t,char*s){
    while(*t!=0)
        t++;

    while((*t++=*s++)!=0)
        ;
}

int main() {
    // char str[] = "Hello";
    // printf("%s\n", str);

    // int b = 108;
    // printf("%c\n", b);

    // char a = 'l';

    // printf("%d\n", a);

    // char o = '0';
    // printf("%d\n", o);

    char d = '!' + '$';
    printf("%c\n", d);

    int p = 'a' + 'A';
    printf("%d\n", p);

    // int g = 'a' + 'B';
    // printf("%d\n", g);

    char ss[10] = "acc", aa[10] = "bbxxyy";
    fun(ss, aa);
    printf("%s,%s\n", ss, aa);

    _Bool h = NULL;
    printf("%d\n", h);

    char cc[] = "abc";
    printf("%s\n", cc);

    double pi = 3.1415926;
    float f = 0.123456789f;
    printf("pi = %f, f = %4f\n", pi, f);
    printf("pi = %d, f = %c\n", pi, f);

    return 0;
}