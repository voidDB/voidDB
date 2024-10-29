#include "include/void.h"

#include <stdio.h>
#include <stdlib.h>

int void_write_zero_file(size_t nBytes, char *fileName)
{
	int e;

	FILE *file = fopen(fileName, "w");
	if (file == NULL) {
		e = 8;

		goto epilogue;
	}

	void *array = calloc(nBytes, 1);

	size_t written = fwrite(array, 1, nBytes, file);
	if (written != nBytes) {
		e = 9;

		goto epilogue;
	}

	e = fclose(file);
	if (e != 0) {
		e = 10;

		goto epilogue;
	}

	file = NULL;

epilogue:
	if (file != NULL) {
		fclose(file);
	}

	free(array);

	return e;
}
