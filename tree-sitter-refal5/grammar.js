/**
 * @file Refal5 grammar for tree-sitter
 * @author Kirill Kiselev <kirillkiselev2003@yandex.ru>
 * @license MIT
 */

/// <reference types="tree-sitter-cli/dsl" />
// @ts-check

module.exports = grammar({
  name: "refal5",
  
  extras: $ => [
    /\s/,
    $.comment
  ],
  
  rules: {
    source_file: $ => repeat(
      choice(
        $.external_declaration,
        $.function_definition,
      )
    ),
    
    external_declaration: $ => seq(
      $.external_modifier,
      field('func_name_list', $.function_name_list),
      ';'
    ),

    function_name_list: $ => repeat1(
      seq(
        field('name', $.ident)
      )
    ),
    
    function_definition: $ => seq( field('entry' ,optional($.entry_modifier)),
      field('name', $.ident,),
      field('body', $.body),
    ),

    body: $ => seq(
      '{',
      repeat1(
        seq($.sentence, ';')
      ),
      optional($.sentence),
      '}'
    ),

    sentence: $ => choice(
      $.sentenceEq, 
      $.sentenceBlock 
    ),

    sentenceEq: $ => seq(
      field(
        "lhs",
        seq(
          optional($._pattern), 
          repeat($.condition)
        )
      ),
      "=",
      field("rhs", optional($._expr)),
    ),

    sentenceBlock: $ => seq(
      optional($._pattern), 
      repeat($.condition),
      ',',
      field('block', $.sentenceBlockEnd),
    ),
    
    sentenceBlockEnd: $ => seq(
      field('expr', optional($._expr)), 
      ':',
      '{',
      field('body', $.body), 
      '}'
    ),
    
    _pattern: $ => repeat1(choice(
      $.ident,
      $.string,
      $.number,
      $.variable,
      $.grouped_expr,
    )),

    grouped_pattern: $ => seq('(', optional($._pattern), ')'),

    condition: $ => seq(
        ',',
        field('result', optional($._expr)),
        ':',
        field('pattern', optional($._pattern)),
    ),

    _expr: $ => repeat1(
      choice(
        $.ident,
        $.string,
        $.number,
        $.variable,
        $.function_call,
        $.grouped_expr,
      )
    ),

    function_call: $ => seq(
      '<',
      field(
        "ident",
        $.ident
      ),
      field(
        "expr",
        optional($._expr)
      ),
      '>'
    ), 

    grouped_expr: $ =>  seq(
      '(',
      optional($._expr),
      ')'
    ),
    
    variable: $ => seq(
      field('type', $.type),
      '.',
      field('name', $.ident),
    ), 
    
    ident : $ => /(([A-Za-z][A-Za-z0-9_-]*)|([0-9]+))/, 

    type: $ => choice(
      's',
      'e',
      't'
    ),

    string: $ => /\"[^\n]*\"/,
    
    number: $ => /('-'|'+')?\d+/,
    
    entry_modifier: $ => '$ENTRY',
    
    external_modifier: $ => choice(
        '$EXTERNAL',
        '$EXTERN',
        '$EXTRN'
    ),

    comment: $ => token(choice(
      seq('//', /[^\n]*/),
      seq('/*', /[^*]*\*+([^/*][^*]*\*+)*/, '/')
    )),

  }
});
