#include <stdio.h>
// #include <ctype.h>
#define SIZE 10
#define BUFSIZE 1000
char buf[BUFSIZE];
int bufp = 0;

void ungetch(int c)
{
	if (bufp >= BUFSIZE)
		printf("ungetch: too many characters\n");
	else
		buf[bufp++] = c;
}

int getch(void)
{
	return (bufp > 0) ? buf[--bufp] : getchar();
}

int main()
{
	int n, c, array[SIZE] = {0}, getint(int *);
	for (n = 0; n < SIZE && (c = getint(&array[n])) != EOF; n++)
	{
		if (c == 0)
		{
			getch();
			n--;
		}
	}
	for (n = 0; n < SIZE; n++)
		printf("%d ", array[n]);
	// int n, array[SIZE] = {0}, getint(int *);
	// for (n = 0; n < SIZE && getint(&array[n]) != EOF; n++) {
	// 	printf("%d ", array[n]);
	// }
}

#include <ctype.h>
																										   int
																										   getch(void);
void ungetch(int);

int getint(int *pn)
{
	int c, sign;
	while (isspace(c = getch()))
		;
	if (!isdigit(c) && c != EOF && c != '+' && c != '-')
	{
		ungetch(c);
		return 0;
	}
	sign = (c == '-') ? -1 : 1;
	if (c == '+' || c == '-')
		c = getch();
	for (*pn = 0; isdigit(c); c = getch())
		*pn = 10 * *pn + c - '0';
	*pn *= sign;
	if (c != EOF)
		ungetch(c);
	return c;
}

