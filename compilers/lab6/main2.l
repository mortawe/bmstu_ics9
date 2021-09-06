%option noyywrap bison-bridge bison-locations

%{
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#define TAG_KEYWORD 1
#define TAG_NUMERICAL_LITERAL 2
#define TAG_IDENT 3
#define TAG_ERROR 4

char* tag_names[] = {"EOF", "KEYWORD", "NUMERICAL_LITERAL", "IDENT", "ERROR"};

typedef struct Position Position;

struct Position {
    int line, pos, index;
};

void print_pos(Position *p) {
    printf("(%d,%d)", p->line, p->pos);
}

struct Fragment {
    Position starting, following;
};

typedef struct Fragment YYLTYPE;
typedef struct Fragment Fragment;

void print_frag(Fragment* f) {
    print_pos(&(f->starting));
    printf("-");
    print_pos(&(f->following));
}

union Token {
    long num_literal;
    char *ident;
    char *keyword;
};

typedef union Token YYSTYPE;

int continued;
struct Position cur;
#define YY_USER_ACTION {                 \
    if (!continued) {                    \
        yylloc->starting = cur;          \
    }                                    \
    continued = 0;                       \
    for (int i = 0; i < yyleng; i++) {   \
        if (yytext[i] == '\n') {         \
            cur.line++;                  \
            cur.pos = 1;                 \
        } else {                         \
            cur.pos++;                   \
	  }                                \
        cur.index++;                     \
    }                                    \
    yylloc->following = cur;             \
}
void err (char *msg){
    printf ("Error");
    print_pos(&cur);
    printf(":%s\n",msg);
}

typedef struct {
    char** array;
    size_t used;
    size_t size;
} Array;

void initArray(Array* a, int initial_size) {
    a->array = malloc(initial_size * sizeof(char*));
    a->used = 0;
    a->size = initial_size;
}

void insertArray(Array* a, char* element) {
    if (a->used == a->size) {
        a->size *= 2;
        a->array = realloc(a->array, a->size * sizeof(char*));
    }
    a->array[a->used++] = element;
}

bool findInArray(Array* a, char* element) {
    for (int i = 0; i < a->used; i++) {
        if (strcmp(a->array[i], element) == 0) {
            return true;
        }
    }
    return false;
}

void freeArray(Array* a) {
    free(a->array);
    a->array = NULL;
    a->used = a->size = 0;
}

Array ids;

void init_scanner(char* program) {
    continued = 0;
    cur.line = 1;
    cur.pos = 1;
    cur.index = 0;
    initArray(&ids, 2);
    yy_scan_string(program);
}




%}

NUMERICAL_LITERAL ([0-9][0-9a-fA-F]*)
IDENT [a-zA-Z][a-zA-Z0-9]*
KEYWORD (qeq)|(xxx)|(xx)
%x SKIP

%%

[\n\t ]+

{NUMERICAL_LITERAL}  {
    long number = (long)strtol(yytext, NULL, 16);
    yylval->num_literal = number;
    return TAG_NUMERICAL_LITERAL;
}
{KEYWORD} {
      yylval->keyword = yytext;
      return TAG_KEYWORD;
}
{IDENT} {
    yylval->ident = yytext;
    if (yytext[strlen(yytext)-1] != yytext[0]) {
        BEGIN(SKIP);
        return TAG_ERROR;
    }
    insertArray(&ids, yytext);
    return TAG_IDENT;
}

. {
    err("unexpected symbol");
    BEGIN(SKIP);
}

<SKIP>[^\n\t ]+

<SKIP>[\n\t ]+ BEGIN(0);

<<EOF>> return 0;

%%




int main() {
    int tag;
    YYSTYPE value;
    YYLTYPE coords;

    FILE *inputfile;
	long size_str;
	char *str;
	union Token token;
	inputfile = fopen("test.txt","r");
	if (inputfile == NULL) {
		fputs("File not found", stderr);
		exit(1);
	}
	fseek(inputfile, 0,SEEK_END);
	size_str = ftell(inputfile);
    rewind(inputfile);

    str=(char*)malloc(sizeof(char)*(size_str+1));
    if (str == NULL) {
		fputs("Memory error",stderr);
		exit(2);
	}
    size_t n = fread(str,sizeof(char),size_str,inputfile);
    if (n != size_str) {
		fputs ("Reading error",stderr);
		exit (3);
	}
    str[size_str] = '\0';
    fclose (inputfile);
    init_scanner(str);
    do{
        tag = yylex(&value,&coords);
        if (tag != 0 && tag != TAG_ERROR) {
            printf("%s ", tag_names[tag]);
            print_frag(&coords);
            if (tag == TAG_NUMERICAL_LITERAL) {
                printf(": %ld", value.num_literal);
		} else if (tag == TAG_IDENT) {
		    printf(": %s", value.ident);
		} else if (tag == TAG_KEYWORD) {
		    printf(": %s", value.keyword);
		}
		printf("\n");
        }
    } while (tag != 0);
    freeArray(&ids);
    return 0;
}
