# CHANGELOG

## v2.1.0

- Update codesystem
- Remove useless metadata ID
- Change FDJID to normalize with the new set of data from the FDJ'API
  - format respect this buid: 'annee_numero_de_tirage'-'1er_ou_2eme_tirage'-'numeric-day-of-month'-'litteral-month'
  - '1er_ou_2eme_tirage' is 1 or 2 in version 1 and always 1 for other version.
- Implement new method to get the planned draw from the FDJ (new API endpoint).
