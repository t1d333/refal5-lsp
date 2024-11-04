import XCTest
import SwiftTreeSitter
import TreeSitterRefal5

final class TreeSitterRefal5Tests: XCTestCase {
    func testCanLoadGrammar() throws {
        let parser = Parser()
        let language = Language(language: tree_sitter_refal5())
        XCTAssertNoThrow(try parser.setLanguage(language),
                         "Error loading Refal5 grammar")
    }
}
