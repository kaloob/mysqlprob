#include <stdio.h> 
#include <stdlib.h>
#include <string.h>

int max(int a, int b){
  if ( a < b ) {
    return b;
  }
  return a;
}

void printer(int **score_matriex, int width, int heigth) {
  int i, j;
  for (i = 0; i < width; i++){
    for (j = 0; j < heigth; j++) {
      printf("%d ", score_matriex[i][j]);
    }
    printf("\n");
  }

}

int main() {
  int i, j;
  char sql1[] = "select * from poe where name = 'poe123'";
  char sql2[] = "select * from poe where name = 'poe123' or 1=1 --'";

  size_t width = strlen(sql1) + 1;
  size_t heigth = strlen(sql2) + 1;

  int **score_matriex = malloc(sizeof(int*) * width);
  for (i = 0; i < width; i++){
    score_matriex[i] = (int *) malloc(sizeof(int) * heigth);
  }

  // initialize
  for (i = 0; i < width; i++){
    for (j = 0; j < heigth; j++) {
      score_matriex[i][j] = 0;
    }
  }

  for (i = 1; i < width; i++){
    for (j = 1; j < heigth; j++) {
      int a = score_matriex[i - 1][j - 1] + sim;
      int b = score_matriex[i - 1][j] - p;
      int c = score_matriex[i][j - 1] - p;
      score_matriex[i][j] = max(max(a, b), max(a, c));
    }
  }

  printer(score_matriex, width, heigth);
}
