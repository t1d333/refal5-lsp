#include "tree_sitter/parser.h"
#include <stdbool.h>

enum TokenType { LINE_COMMENT };

void *tree_sitter_refal5_external_scanner_create(void) { return NULL; }

void tree_sitter_refal5_external_scanner_destroy(void *payload) {}

unsigned tree_sitter_refal5_external_scanner_serialize(void *payload,
                                                       char *buffer) {
  return 0;
}

void tree_sitter_refal5_external_scanner_deserialize(void *payload,
                                                     const char *buffer,
                                                     unsigned length) {}

bool tree_sitter_refal5_external_scanner_scan(void *payload, TSLexer *lexer,
                                              const bool *valid_symbols) {

  while ((lexer->lookahead == '\n') || (lexer->lookahead == '\r')) {
    lexer->advance(lexer, false);
  }

  if (lexer->lookahead == '*' && lexer->get_column(lexer) == 0) {

    while (lexer->lookahead != '\n' && !lexer->eof(lexer)) {
      lexer->advance(lexer, false);
    }

    lexer->result_symbol = LINE_COMMENT;
    return true;
  }

  return false;
}
