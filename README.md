# GeWoScout 

## TL;DR
This is an application which helps people find their dream Genossenschaftswohnung in Vienna. It is a service which consolidates the existing offerings of 3 biggest Genossenschaften in Vienna and allows you to enlist for your chosen appartment in no time!

## Target audience
People who are looking for an affordable place to stay in Vienna.

## What exactly are Genossenschaftswohnungen
Genossenschaftswohnungen are appartments that are rented directly by the company (Genossenschaft) who built the building (or currently owns it) and usually cost less than the standard apartments for rent. The catch is that you have to pay an Affiliation Fee (Genossenschaftsanteil) which can sometimes be very high (it is about 15 000€ on average). However, even though you might have a lot of costs to cover at the beginning of renting, those apartments are in very high demand. When you move out, you get the most of your Affiliation Fee back and the rent contract you sign with the company **does not have an end date**. This basically means that you can live in such an apartment as if it were your own 😎. 

## How does our application fit in
Because those apartments are so **juicy**, they usually *sell* like hot cakes. Unfortunately, every Genossenschaft has their own website, where they post their currently available apartments. And because there is no official and free platform which would list all the available apartments as soon as they appear on the websites of the Genossenschaften, the hunt for such an appartment becomes extremely difficult.

We would like to propose a service that allows you to see all available apartments in one place without having to switch to 10 different websites every day. As an additional helper, a notification bot will send you new apartments as soon as they appear.

## Genossenschaften that we are supporting
- BWS Gruppe https://www.bwsg.at/
- Österreichisches Volkswohnungswerk (ÖVW) https://www.oevw.at/
- WBV-GPA https://www.wbv-gpa.at/

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

## Notes for the reviewer
The application components are not packaged in containers, so there is no Dockerfile in the repo (and  there are also no artifacts). CI testing is done using GitHub actions workflows (please find the appropriate files in the `.github/worflows` directory). There are pipelines for changes to our frontend, backend and infrastructure components that test the changes on every PR. The backend CI contains integration testing. There is no versioning scheme for the components. The project can be entirely deployed and destroyed with our `tf_apply` pipeline (the required storage accounts and a Static Web Apps instance have to present beforehand). The `infra` directory contains terraform configuration of our infrastructure and the `README.md` in that directory has short instructions how to run everything. More details (such as our Decision Records) can be found in the project's Wiki.
