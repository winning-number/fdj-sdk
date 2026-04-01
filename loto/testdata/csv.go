package testdata

var (
	InvalidCSV = []byte(
		`toto,tata,truc
41,42,44`)

	ValidCSVV0 = []byte(
		`annee_numero_de_tirage;jour_de_tirage;date_de_tirage;date_de_forclusion;boule_1;boule_2;boule_3;boule_4;boule_5;boule_6;boule_complementaire;combinaison_gagnante_en_ordre_croissant;nombre_de_gagnant_au_rang1;rapport_du_rang1;nombre_de_gagnant_au_rang2;rapport_du_rang2;nombre_de_gagnant_au_rang3;rapport_du_rang3;nombre_de_gagnant_au_rang4;rapport_du_rang4;nombre_de_gagnant_au_rang5;rapport_du_rang5;nombre_de_gagnant_au_rang6;rapport_du_rang6;nombre_de_gagnant_au_rang7;rapport_du_rang7;numero_jokerplus;devise;
2008002;VE;20080613;20080813;34;25;37;42;15;36;49;15-25-34-36-37-42;2;7501991;6;78102,6;225;5101,4;613;242,4;14049;121,2;17651;16,2;253160;8,1;8 651 700;eur;
2008001;VE;20080229;20080430;1;49;30;17;15;13;48;1-13-15-17-30-49;2;10001036;8;44953,9;351;2493,7;574;161,8;16166;80,9;15033;12,2;259505;6,1;8 152 912;frf;
`)
	ValidCSVV0WithWarning = []byte(
		`annee_numero_de_tirage;jour_de_tirage;date_de_tirage;date_de_forclusion;boule_1;boule_2;boule_3;boule_4;boule_5;boule_6;boule_complementaire;combinaison_gagnante_en_ordre_croissant;nombre_de_gagnant_au_rang1;rapport_du_rang1;nombre_de_gagnant_au_rang2;rapport_du_rang2;nombre_de_gagnant_au_rang3;rapport_du_rang3;nombre_de_gagnant_au_rang4;rapport_du_rang4;nombre_de_gagnant_au_rang5;rapport_du_rang5;nombre_de_gagnant_au_rang6;rapport_du_rang6;nombre_de_gagnant_au_rang7;rapport_du_rang7;numero_jokerplus;devise;truc;
2008002;VE;20080613;20080813;34;25;37;42;15;36;49;15-25-34-36-37-42;2;7501991;6;78102,6;225;5101,4;613;242,4;14049;121,2;17651;16,2;253160;8,1;8 651 700;eur;42;
2008001;VE;20080229;20080430;1;49;30;17;15;13;48;1-13-15-17-30-49;2;10001036;8;44953,9;351;2493,7;574;161,8;16166;80,9;15033;12,2;259505;6,1;8 152 912;frf;21;
`)
	ValidCSVHeader = []byte(
		`annee_numero_de_tirage;jour_de_tirage;date_de_tirage;date_de_forclusion;boule_1;boule_2;boule_3;boule_4;boule_5;boule_6;boule_complementaire;combinaison_gagnante_en_ordre_croissant;nombre_de_gagnant_au_rang1;rapport_du_rang1;nombre_de_gagnant_au_rang2;rapport_du_rang2;nombre_de_gagnant_au_rang3;rapport_du_rang3;nombre_de_gagnant_au_rang4;rapport_du_rang4;nombre_de_gagnant_au_rang5;rapport_du_rang5;nombre_de_gagnant_au_rang6;rapport_du_rang6;nombre_de_gagnant_au_rang7;rapport_du_rang7;numero_jokerplus;devise;
`)
)
