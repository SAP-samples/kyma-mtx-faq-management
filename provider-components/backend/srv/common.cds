/*
	Common Annotations shared by all apps
*/

using { sap.demo.faq as my } from '../db/schema.cds';


////////////////////////////////////////////////////////////////////////////
//
//	FAQ Lists
//
annotate my.Faqs with @(
	Common.SemanticKey: [ID],
	UI: {
		Identification: [{Value:ID}],
		SelectionFields: [ ID, author_ID, count ],
		LineItem: [
			{Value: ID},
			{Value: title, Label:'{i18n>Question}'},
			{Value: author.name, Label:'{i18n>Author}'},
			{Value: descr, Label:'{i18n>Description}'},
			{Value: category.name},
			{Value: state},
			{Value: count}
		]
	}
) {
	author @ValueList.entity:'Authors';
};

////////////////////////////////////////////////////////////////////////////
//
//	FAQs Details
//
annotate my.Faqs with @(
	UI: {
		HeaderInfo: {
			TypeName: '{i18n>Faq}',
			TypeNamePlural: '{i18n>Faqs}',
			Title: {Value: title},
			Description: {Value: descr}
		},
	}
);

////////////////////////////////////////////////////////////////////////////
//
//	Faqs Elements
//
annotate my.Faqs with {
	ID @title:'{i18n>ID}' @UI.HiddenFilter;
	title @title:'{i18n>Title}';
	category  @title:'{i18n>Category}'  @Common: { Text: category.name,  TextArrangement: #TextOnly };
	author @title:'{i18n>Author}' @Common: { Text: author.name, TextArrangement: #TextOnly };
	count @title:'{i18n>Count}';
	state @title:'{i18n>State}';
	descr @title:'{i18n>Description}';
	descr @UI.MultiLineText;
	answer @UI.MultiLineText;
}

////////////////////////////////////////////////////////////////////////////
//
//	Categories Elements
//
annotate my.Categories with {
	ID  @title: '{i18n>ID}';
	name  @title: '{i18n>Category}';
	parent  @title: '{i18n>Parent}';
}

////////////////////////////////////////////////////////////////////////////
//
//	Categories 
//
annotate my.Categories with @(
	Common.SemanticKey: [name],
	UI: {
		SelectionFields: [ name ],
		Identification: [{Value:name}],
		HeaderInfo: {
			TypeName: '{i18n>Category}',
			TypeNamePlural: '{i18n>Categories}',
			Title: {Value: name},
			Description: {Value: parent_ID}
		},
		LineItem:[
			{Value: ID},
			{Value: name},
			{Value: parent.name, Label: 'Main Category'},
		],
	}
);

////////////////////////////////////////////////////////////////////////////
//
//	Authors Elements
//
annotate my.Authors with {
	ID @title:'{i18n>ID}' @UI.HiddenFilter;
	name @title:'{i18n>Name}';
}

//
//	Authors 
//
annotate my.Authors with @(
	Common.SemanticKey: [name],
	UI: {
		Identification: [{Value:name}],
		SelectionFields: [ name ],
		HeaderInfo: {
			TypeName: '{i18n>Author}',
			TypeNamePlural: '{i18n>Authors}',
			Title: {Value: name},
		},
		LineItem:[
			{Value: ID},
			{Value: name},
		],
	}
);

////////////////////////////////////////////////////////////////////////////
//
//	Languages List
//
annotate common.Languages with @(
	Common.SemanticKey: [code],
	Identification: [{Value:code}],
	UI: {
		SelectionFields: [ name, descr ],
		LineItem:[
			{Value: code},
			{Value: name},
		],
	}
);

////////////////////////////////////////////////////////////////////////////
//
//	Language Details
//
annotate common.Languages with @(
	UI: {
		HeaderInfo: {
			TypeName: '{i18n>Language}',
			TypeNamePlural: '{i18n>Languages}',
			Title: {Value: name},
			Description: {Value: descr}
		},
		Facets: [
			{$Type: 'UI.ReferenceFacet', Label: '{i18n>Details}', Target: '@UI.FieldGroup#Details'},
		],
		FieldGroup#Details: {
			Data: [
				{Value: code},
				{Value: name},
				{Value: descr}
			]
		},
	}
);
