# LSB Go

## Principe
Ce programme en GO se base sur la méthode LSB pour cacher des informations.
Pour cela il faut donner une image en format PNG en entrée et un fichier à cacher.
Le message sera chiffré puis inseré dans l'image suivant un certains parcours.
Pour l'extraction, l'opération inverse est appliqué pour récuperer les données puis une fonction de déchiffrement est utilisé.
Le chiffrement est basic mais il utilise une clé de taille variable.
Une entête est utilisée pour garder la taille des données à extraire.

## Les différents parcours

Ce programme utilise différentes méthodes de parcours de l'image :
- horizontal : les pixels sont parcourus ligne par ligne.
- vertical :  les pixels sont parcourus colonne par colonne.
- diagonal : les pixels sont parcourus de manière diagonal
- corps fini : les pixels sont parcourus de manière à suivre la génération d'un corps fini. Cela permet d'avoir 2 paramètres supplémentaires qui sont l'ordre du corps et un générateur de ce corps.


## Chiffrement

Le chiffrement se fait avec un simple xor et suit la méthodologie du CBC.
Le vecteur d'initialisation est modifié à chaque tour.
La clé est de taille variable, le chiffrement se fait par bloc d'un octet.
De ce fait on parcours les octets de le clé en même temps que le message à un modulo près.




