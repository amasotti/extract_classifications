# Classify (K10+)

**Work in progress**


The idea is to create a small flexible tool to extract classifications (rvk, bkl, ddc, lcc) from the GVK Catalogue

+ [Sru](https://wiki.k10plus.de/display/K10PLUS/SRU) address

      https://sru.k10plus.de/gvk


The SRU protocol supports several output formats. This scritp works (at moment at least) only with the [MODS](https://en.wikipedia.org/wiki/Metadata_Object_Description_Schema) format.


## Temporary: Testing

I'm currently testing the XMLParser. At moment the script doesn't use the HTTP pkg to get the xml but uses a xml file obtained with the following request:
      https://sru.k10plus.de/gvk?version=1.1&operation=searchRetrieve&query=pica.thm=Dostoevsky&recordSchema=mods&maximumRecords=100

## Subjects Headings

### Authority:

+ **bkl** : [Basisklassifikation](https://www.gbv.de/bibliotheken/verbundbibliotheken/02Verbund/01Erschliessung/02Richtlinien/05Basisklassifikation/index)
+ **bisacsh** : [BISAC Subject Headings](https://bisg.org/page/bisacedition)
+ **bicss** : [BIC Subject Categories](https://bic.org.uk/files/pdfs/101201%20bic2.1%20complete%20rev.pdf)
+ **ddc** : Dewey Class or Division
+ **fid** : Fachinformationsdienst
+ **lcsh** : [Library of Congress Subject Headings](https://id.loc.gov/vocabulary/subjectSchemes/bisacsh.html)
+ **rvk** : [Regensburger Verbundsklassifikation](https://rvk.uni-regensburg.de/regensburger-verbundklassifikation-online)

Check [loc.gov](https://www.loc.gov/standards/sourcelist/subject.html) for a complete list of authority codes.


### Example of an output in xml format:
(<small>Formatted for readability</small>)
```xml
    <zs:searchRetrieveResponse>
      <zs:version>1.1</zs:version>
      <zs:numberOfRecords>1007</zs:numberOfRecords>
      <zs:records>
        <zs:record>
          <zs:recordSchema>mods</zs:recordSchema>
          <zs:recordPacking>xml</zs:recordPacking>
          <zs:recordData>
            <mods version="3.6" xsi:schemaLocation="http://www.loc.gov/mods/v3 http://www.loc.gov/standards/mods/v3/mods-3-6.xsd">

              <titleInfo>
                <title>Irène Némirovsky's Russian Influences</title>
                <subTitle>Tolstoy, Dostoevsky and Chekhov</subTitle>
              </titleInfo>

              <name type="personal" usage="primary">
                <namePart>Cenedese, Marta-Laura</namePart>
                  <role>
                    <roleTerm type="text">VerfasserIn</roleTerm>
                  </role>
                  <role>
                    <roleTerm authority="marcrelator" type="code">aut</roleTerm>
                </role>
              </name>
              <typeOfResource>text</typeOfResource>
              <genre authority="rdacontent">Text</genre>

              <originInfo>
                <place>
                  <placeTerm type="code" authority="marccountry">gw</placeTerm>
                </place>
                <place>
                  <placeTerm type="code" authority="iso3166">XA-DE</placeTerm>
                </place>
                <dateIssued encoding="marc">2021</dateIssued>
                <edition>1st ed. 2021.</edition>
                <issuance>monographic</issuance>
                </originInfo>
                <originInfo eventType="publication">
                  <place></place>
                  <publisher>Springer International Publishing</publisher>
                  <dateIssued>2021.</dateIssued>
              </originInfo>
              
              <originInfo eventType="publication">
                <place>
                  <placeTerm type="text">Cham</placeTerm>
                </place>
                <publisher>Imprint: Palgrave Macmillan</publisher>
                <dateIssued>2021.</dateIssued>
              </originInfo>
              <language>
                <languageTerm authority="iso639-2b" type="code">eng</languageTerm>
              </language>
              
              <physicalDescription>
                <form authority="marccategory">electronic resource</form>
                <form authority="marcsmd">remote</form>
                <extent>1 Online-Ressource(XX, 213 p. 1 illus.)</extent>
                <form type="media" authority="rdamedia">Computermedien</form>
                <form type="carrier" authority="rdacarrier">Online-Ressource</form>
              </physicalDescription>
              
              <abstract type="Summary">
              1. Introduction; Creative Encounters over Time and Space: Writers, Readers, and Researchers -- 2. From Russia to France, via England: Suite française, War and Peace, and E. M. Forster -- 3. Departing from Tolstoy: Polyphony and Monologism -- 4. Beyond Tolstoy: Music -- 5. Dreams from Underground -- 5. The Abject -- 6. An Anthropology of Suffering -- 7. La Vie de Tchekhov: A Romanced Biography -- 8. La Vie de Tchekhov in the 21st Century -- 9. Conclusion: A Russian Suite.
              </abstract>

              <abstract type="Summary">
              This book explores the influence of Tolstoy, Dostoevsky, and Chekhov on Russian-born French language writer Irène Némirovsky. It considers the complexity of each of these relationships and the different modes in which they appear; demonstrating how, by skillfully integrating reading and writing, reception and creation, Némirovsky engaged with Russian literature within her own work. Through detailed analysis of the intersections between novels, short stories and archival sources, the book assesses to what degree Tolstoy, Dostoevsky and Chekhov influenced Némirovsky, how this influence affected her work, and to what effects. To this aim the book articulates the notion of creative influence, a method that, in conversation with theories of influence, intertextuality, and reception aesthetics, seeks to reflect a “meeting of artistic minds” that includes affective, ethical, and creative encounters between writers, readers, and researchers.
              </abstract>

              <note type="statement of responsibility" altRepGroup="00">by Marta-Laura Cenedese</note>
              <subject authority="lcsh">
                <topic>Literature, Modern—20th century</topic>
              </subject>
                <subject authority="lcsh">
              <topic>European literature</topic>
                </subject>
              <subject authority="lcsh">
                <topic>Comparative literature</topic>
              </subject>
              
              <classification authority="ddc" edition="23">809.04</classification>
              <classification authority="bicssc">DSBH</classification>
              <classification authority="bisacsh">LIT024050</classification>
              
              <location>
                <url displayLabel="electronic resource" usage="primary display" note="Lizenzpflichtig">https://doi.org/10.1007/978-3-030-44203-3</url>
              </location>
              
              <relatedItem type="series">
                <titleInfo>
                  <title>Palgrave Studies in Modern European Literature</title>
                </titleInfo>
              </relatedItem>
              
              <relatedItem type="series">
                <titleInfo>
                  <title>Springer eBook Collection</title>
                </titleInfo>
              </relatedItem>

              <relatedItem type="otherFormat"/>
              <relatedItem type="otherFormat"/>
              <relatedItem type="otherFormat"/>
              
              <relatedItem type="otherFormat" otherType="Erscheint auch als" displayLabel="Erscheint auch als">
                <note>Druck-Ausgabe</note>
              </relatedItem>
              
              <relatedItem type="otherFormat" otherType="Erscheint auch als" displayLabel="Erscheint auch als">
                <note>Druck-Ausgabe</note>
              </relatedItem>
              
              <relatedItem type="otherFormat" otherType="Erscheint auch als" displayLabel="Erscheint auch als">
                <note>Druck-Ausgabe</note>
              </relatedItem>
              
              <identifier type="isbn">9783030442033</identifier>
              <identifier type="doi">10.1007/978-3-030-44203-3</identifier>
              
              <recordInfo>
                <descriptionStandard>rda</descriptionStandard>
                <recordContentSource authority="marcorg">DE-627</recordContentSource>
                <recordCreationDate encoding="marc">201201</recordCreationDate>
                <recordChangeDate encoding="iso8601">20201201105043.0</recordChangeDate>
                <recordIdentifier source="DE-627">1741583977</recordIdentifier>
              
              <recordOrigin>
                Converted from MARCXML to MODS version 3.6 using MARC21slim2MODS3-6.xsl (Revision 1.119 2018/06/21)
              </recordOrigin>
              
              <languageOfCataloging>
                <languageTerm authority="iso639-2b" type="code">ger</languageTerm>
              </languageOfCataloging>
              
              </recordInfo>
                </mods>
            </zs:recordData>
    <zs:recordPosition>1</zs:recordPosition>
    </zs:record>
```

