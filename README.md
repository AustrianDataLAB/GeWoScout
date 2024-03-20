# GeWoScout 

## TL;DR
This is an application which helps people find their dream Genossenschaftswohnung in Vienna. It is a service which consolidates the existing offerings of 10 biggest Genossenschaften in Vienna and allows you to enlist for your chosen appartment in no time!

## What exactly are Genossenschaftswohnungen
Genossenschaftswohnungen are appartments that are rented directly by the company (Genossenschaft) who built the building (or currently owns it) and usually cost less than the standard apartments for rent. The catch is that you have to pay an Affiliation Fee (Genossenschaftsanteil) which can sometimes be very high (it is about 15 000€ on average). However, even though you might have a lot of costs to cover at the beginning of renting, those apartments are in very high demand. When you move out, you get the most of your Affiliation Fee back and the rent contract you sign with the company **does not have an end date**. This basically means that you can live in such an apartment as if it were your own 😎. 

## How does our application fit in
Because those apartments are so **juicy**, they usually *sell* like hot cakes. Unfortunately, every Genossenschaft has their own website, where they post their currently available apartments. And because there is no official and free platform which would list all the available apartments as soon as they appear on the websites of the Genossenschaften, the hunt for such an appartment becomes extremely difficult.

We would like to propose a service that allows you to see all available apartments in one place without having to switch to 10 different websites every day. As an additional helper, a notification bot will send you new apartments as soon as they appear.

## Genossenschaften that we are supporting
- BWS Gruppe https://www.bwsg.at/
- Österreichisches Volkswohnungswerk (ÖVW) https://www.oevw.at/
- Österreichisches Siedlungswerk (ÖSW) https://www.oesw.at/
- Sozialbau https://www.sozialbau.at/
- Siedlungsunion https://www.siedlungsunion.at/
- Neue Heimat Gewog https://www.nhg.at/
- Familienwohnbau https://familienwohnbau.at/
- Neues Leben https://www.wohnen.at/
- WBV-GPA https://www.wbv-gpa.at/
- Wien Süd https://www.wiensued.at/

## What data do we scrape
We only scrape the most relevant information a user would filter by. To allow users to get more detailed information we provide a link to the details page.
- Title
- Genossenschaft Name
- Location
- Number or rooms
- Area
- Move in date
- Construction year
- Enerey-efficiency class HWG (A++, A+, A, B, C, D, E, F, G)
- Enerey-efficiency class FGEE (A++, A+, A, B, C, D, E, F, G)
- Price monthly
- Price Genossenschafts-Deposit
- URL to Gennossenschafts site details page
