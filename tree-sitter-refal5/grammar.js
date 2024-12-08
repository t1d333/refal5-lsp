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
    $.comment,
    $.line_comment,
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

    function_name_list: $ => 
      choice(
        field('name', $.ident),
        field('name', seq($.ident, repeat1(seq(',', $.ident)))),
      ),
    
    function_definition: $ => seq( field('entry' ,optional($.entry_modifier)),
      field('name', $.ident,),
      field('body', $.body),
      optional(';')
    ),

    body: $ => seq(
      '{',
      repeat(
        seq($.sentence, ';')
      ),
      choice(
        $.sentence,
        seq($.sentence, ';')
      ),
      '}'
    ),

    sentence: $ => choice(
      field("sentence_eq" ,$.sentence_eq), 
      field("sentence_block", $.sentence_block), 
    ),

    sentence_eq: $ => seq(
      optional(field("lhs", $.sentence_lhs)), 
      repeat(field("condition", $.condition)),
      "=",
      optional(field("rhs", $.sentence_rhs)),
    ),

    sentence_block: $ => seq(
      optional(field("lhs", $.sentence_lhs)), 
      repeat(field("condition", $.condition)),
      ',',
      field('block', $.sentence_block_end),
    ),

    sentence_lhs: $ => seq(
      $._pattern
    ),
    
    sentence_rhs: $ => seq(
      $._expr
    ),
    
    sentence_block_end: $ => seq(
      field('expr', optional($._expr)), 
      ':',
      field('body', $.body), 
    ),
    
    _pattern: $ => repeat1(choice(
      $.ident,
      $.string,
      $.number,
      $.variable,
      $.grouped_pattern,
      $.symbols
    )),

    grouped_pattern: $ => seq('(', optional($._pattern), ')'),

    condition: $ => seq(
        ',',
        optional(field('result', $._expr)),
        ':',
        optional(field('pattern', $._pattern)),
    ),

    _expr: $ => repeat1(
      choice(
        $.ident,
        $.string,
        $.number,
        $.variable,
        $.function_call,
        $.grouped_expr,
        $.symbols
      )
    ),

    function_call: $ => seq(
      '<',
      
      field(
        "name",
        choice(
          $.ident,
          '+',
          '-',
          '*',
          '/',
        )
      ),
      field(
        "param",
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
    
    ident : $ => token(/(([A-Za-z][A-Za-z0-9_-]*)|([0-9]+))/), 

    type: $ => choice(
      's',
      'e',
      't'
    ),

    string: $ => seq(
      /"/,
      repeat(choice(
        prec(10, repeat1("\\\\")),
        prec(-10, '\\"'),
        /[^"\n]/,
      )),
      /"/,
    ),
    
    number: $ => /('-'|'+')?\d+/,
    
    symbols: $ => seq(
      /\'/,
      repeat(choice(
        prec(10, repeat1("\\\\")),
        prec(-10, "\\'"),
        /[^'\n]/
      )),
      /\'/,
    ),

    entry_modifier: $ => token(prec(10,'$ENTRY')),
    
    external_modifier: $ => token(prec(10, choice(
        '$EXTERNAL',
        '$EXTERN',
        '$EXTRN'
    ))),

    comment: $ => token(prec(-1, seq(
      '/*',
      repeat(choice(
        /[^\*]/,
        seq('*', /[^/]/)
      )),
      '*/'             
    ))),

    line_comment: $ => token(prec(-1, seq(
      '*',               
      /[^\n]*/          
    ))),
  }
});
