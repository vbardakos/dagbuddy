package lsp

import "encoding/json"

type InsertMode int
type InsertFormat int
type CompletionItemKind int

const (
	AsIs InsertMode = iota + 1
	AdjustIndendation

	/**
	 * The primary text to be inserted is treated as a snippet.
	 *
	 * A snippet can define tab stops and placeholders with `$1`, `$2`
	 * and `${3:foo}`. `$0` defines the final tab stop, it defaults to
	 * the end of the snippet. Placeholders with equal identifiers are linked,
	 * that is typing in one will update others too.
	 */
	PlainText InsertFormat = iota + 1
	Snippet

	CompletionText CompletionItemKind = iota + 1
	CompletionMethod
	CompletionFunction
	CompletionConstructor
	CompletionField
	CompletionVariable
	CompletionClass
	CompletionInterface
	CompletionModule
	CompletionProperty
	CompletionUnit
	CompletionValue
	CompletionEnum
	CompletionKeyword
	CompletionSnippet
	CompletionColor
	CompletionFile
	CompletionReference
	CompletionFolder
	CompletionEnumMember
	CompletionConstant
	CompletionStruct
	CompletionEvent
	CompletionOperator
	CompletionTypeParameter
)

/**
 * Represents a collection of [completion items](#CompletionItem) to be
 * presented in the editor.
 */
type CompletionList struct {
	IsIncomplete bool                    `json:"isIncomplete"`
	ItemDefaults *CompletionItemDefaults `json:"itemDefaults"`
	Items        []CompletionItem        `json:"items"`
}

type CompletionItem struct {
	Label      string             `json:"label"`
	Kind       CompletionItemKind `json:"kind"`
	InsertText string             `json:"insertText,omitempty"`
	InsertMode InsertMode         `json:"insertTextMode"`
	Deprecated bool               `json:"deprecated"`
	TextEdit   *TextEdit          `json:"textEdit"`
}

/**
  - An edit which is applied to a document when selecting this completion.
  - When an edit is provided the value of `insertText` is ignored.
    *
  - *Note:* The range of the edit must be a single line range and it must
  - contain the position at which completion has been requested.
    *
  - Most editors support two different operations when accepting a completion
  - item. One is to insert a completion text and the other is to replace an
  - existing text with a completion text. Since this can usually not be
  - predetermined by a server it can report both ranges. Clients need to
  - signal support for `InsertReplaceEdit`s via the
  - `textDocument.completion.completionItem.insertReplaceSupport` client
  - capability property.
    *
  - *Note 1:* The text edit's range as well as both ranges from an insert
  - replace edit must be a [single line] and they must contain the position
  - at which completion has been requested.
  - *Note 2:* If an `InsertReplaceEdit` is returned the edit's insert range
  - must be a prefix of the edit's replace range, that means it must be
  - contained and starting at the same position.
    *
  - @since 3.16.0 additional type `InsertReplaceEdit`
*/
// TextEdit | InsertReplaceEdit;
type TextEdit struct {
	Range Range  `json:"range"`
	Text  string `json:"newText"`
}

/**
 * In many cases the items of an actual completion result share the same
 * value for properties like `commitCharacters` or the range of a text
 * edit. A completion list can therefore define item defaults which will
 * be used if a completion item itself doesn't specify the value.
 *
 * If a completion list specifies a default value and a completion item
 * also specifies a corresponding value the one from the item is used.
 *
 * Servers are only allowed to return default values if the client
 * signals support for this via the `completionList.itemDefaults`
 * capability.
 *
 * @since 3.17.0
 */
type CompletionItemDefaults struct {
	CommitChars     []string             `json:"commitCharacters,omitempty"`
	EditRange       *CompletionEditRange `json:"editRange,omitempty"`
	InsertTxtFormat InsertFormat         `json:"insertTextFormat,omitempty"`
	InsertTxtMode   InsertMode           `json:"insertTextMode,omitempty"`
	Data            json.RawMessage      `json:"data,omitempty"`
}

type CompletionEditRange struct {
	Insert  Range `json:"insert"`
	Replace Range `json:"replace"`
}

// export interface CompletionItem {
// 	/**
// 	 * Additional details for the label
// 	 *
// 	 * @since 3.17.0
// 	 */
// 	labelDetails?: CompletionItemLabelDetails;
//
// 	/**
// 	 * Tags for this completion item.
// 	 *
// 	 * @since 3.15.0
// 	 */
// 	tags?: CompletionItemTag[];
//
// 	/**
// 	 * A human-readable string with additional information
// 	 * about this item, like type or symbol information.
// 	 */
// 	detail?: string;
//
// 	/**
// 	 * A human-readable string that represents a doc-comment.
// 	 */
// 	documentation?: string | MarkupContent;
//
//
// 	/**
// 	 * Select this item when showing.
// 	 *
// 	 * *Note* that only one completion item can be selected and that the
// 	 * tool / client decides which item that is. The rule is that the *first*
// 	 * item of those that match best is selected.
// 	 */
// 	preselect?: boolean;
//
// 	/**
// 	 * A string that should be used when comparing this item
// 	 * with other items. When omitted the label is used
// 	 * as the sort text for this item.
// 	 */
// 	sortText?: string;
//
// 	/**
// 	 * A string that should be used when filtering a set of
// 	 * completion items. When omitted the label is used as the
// 	 * filter text for this item.
// 	 */
// 	filterText?: string;
//
// 	/**
// 	 * A string that should be inserted into a document when selecting
// 	 * this completion. When omitted the label is used as the insert text
// 	 * for this item.
// 	 *
// 	 * The `insertText` is subject to interpretation by the client side.
// 	 * Some tools might not take the string literally. For example
// 	 * VS Code when code complete is requested in this example
// 	 * `con<cursor position>` and a completion item with an `insertText` of
// 	 * `console` is provided it will only insert `sole`. Therefore it is
// 	 * recommended to use `textEdit` instead since it avoids additional client
// 	 * side interpretation.
// 	 */
// 	insertText?: string;
//
// 	/**
// 	 * The format of the insert text. The format applies to both the
// 	 * `insertText` property and the `newText` property of a provided
// 	 * `textEdit`. If omitted defaults to `InsertTextFormat.PlainText`.
// 	 *
// 	 * Please note that the insertTextFormat doesn't apply to
// 	 * `additionalTextEdits`.
// 	 */
// 	insertTextFormat?: InsertTextFormat;
//
// 	/**
// 	 * How whitespace and indentation is handled during completion
// 	 * item insertion. If not provided the client's default value depends on
// 	 * the `textDocument.completion.insertTextMode` client capability.
// 	 *
// 	 * @since 3.16.0
// 	 * @since 3.17.0 - support for `textDocument.completion.insertTextMode`
// 	 */
// 	insertTextMode?: InsertTextMode;
//
//
// 	/**
// 	 * The edit text used if the completion item is part of a CompletionList and
// 	 * CompletionList defines an item default for the text edit range.
// 	 *
// 	 * Clients will only honor this property if they opt into completion list
// 	 * item defaults using the capability `completionList.itemDefaults`.
// 	 *
// 	 * If not provided and a list's default range is provided the label
// 	 * property is used as a text.
// 	 *
// 	 * @since 3.17.0
// 	 */
// 	textEditText?: string;
//
// 	/**
// 	 * An optional array of additional text edits that are applied when
// 	 * selecting this completion. Edits must not overlap (including the same
// 	 * insert position) with the main edit nor with themselves.
// 	 *
// 	 * Additional text edits should be used to change text unrelated to the
// 	 * current cursor position (for example adding an import statement at the
// 	 * top of the file if the completion item will insert an unqualified type).
// 	 */
// 	additionalTextEdits?: TextEdit[];
//
// 	/**
// 	 * An optional set of characters that when pressed while this completion is
// 	 * active will accept it first and then type that character. *Note* that all
// 	 * commit characters should have `length=1` and that superfluous characters
// 	 * will be ignored.
// 	 */
// 	commitCharacters?: string[];
//
// 	/**
// 	 * An optional command that is executed *after* inserting this completion.
// 	 * *Note* that additional modifications to the current document should be
// 	 * described with the additionalTextEdits-property.
// 	 */
// 	command?: Command;
//
// 	/**
// 	 * A data entry field that is preserved on a completion item between
// 	 * a completion and a completion resolve request.
// 	 */
// 	data?: LSPAny;
// }
