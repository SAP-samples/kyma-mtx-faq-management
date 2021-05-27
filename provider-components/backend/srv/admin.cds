using AdminService from './admin-service';

////////////////////////////////////////////////////////////////////////////
//
//	Faqs Object Page
//

annotate AdminService.Faqs with {
  ID       @Core.Computed;
  descr    @mandatory;
  title    @mandatory;
  author   @mandatory;
  category @mandatory;
}

annotate AdminService.Authors with {
  ID   @Core.Computed;
  name @mandatory;
}

annotate AdminService.Category with {
  name @mandatory;
}

annotate AdminService.Faqs with @(UI : {
  Facets              : [
  {
    $Type  : 'UI.ReferenceFacet',
    Label  : '{i18n>General}',
    Target : '@UI.FieldGroup#General'
  }
  ],
  FieldGroup #General : {Data : [
  {Value : author_ID},
  {Value : category_ID},
  {Value : count},
  {Value : state}
  ]}
});


////////////////////////////////////////////////////////////
//
//  Draft for Localized Data
//

//annotate AdminService.Faqs with @fiori.draft.enabled;
//annotate AdminService.Faqs with @odata.draft.enabled;

annotate AdminService.Faqs_texts with @(UI : {
  Identification  : [{Value : title}],
  SelectionFields : [
  locale,
  title
  ],
  LineItem        : [
  {
    Value : locale,
    Label : 'Locale'
  },
  {
    Value : title,
    Label : 'Title'
  },
  {
    Value : descr,
    Label : 'Description'
  },
  ]
});
