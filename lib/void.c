#include "include/void.h"

#include <stdio.h>
#include <stdlib.h>

int write_zero_file(size_t nBytes, char *fileName)
{
	FILE *file = fopen(fileName, "w");
	if (file == NULL) {
		return 8;
	}

	void *array = calloc(nBytes, 1);

	size_t written = fwrite(array, 1, nBytes, file);
	if (written != nBytes) {
		return 9;
	}

	if (fclose(file) != 0) {
		return 10;
	}

	return 0;
}
