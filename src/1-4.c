#include <stdio.h>
/* 编写一个程序打印摄氏度转换为相应华氏度的转换表*/

int celsius2fahr(int degree) {
    return degree * 9 / 5 + 32;
}

int printFahrList(int start, int end, int step, int withHead) {
    if (start <= end && step < 0 || start >= end && step > 0 || step == 0) {
        printf("params error !\n");
        return 1;
    }
    if (withHead) {
        printf("celsius\tfahr\n");
    }

    int i = start;
    while (start <= end ? i <= end : i >= end) {
        printf("%d\t%d\n", i, celsius2fahr(i));
        i += step;
    }

    return 0;
}

int main() {
    printFahrList(-17, 148, 11, 1);
}