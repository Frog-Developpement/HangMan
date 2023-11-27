HANGMAN WEB
Projet Hangman Web notice

Le jeu "Hangman Web Game" est une adaptation du célèbre jeu du pendu en version web. Le jeu consiste à deviner un mot en proposant des lettres une par une. 
Chaque mauvaise proposition fait avancer le dessin du pendu, et si le joueur ne parvient pas à deviner le mot avant que le pendu ne soit complet, il perd la partie.

L'objectif du jeu est de mettre à l'épreuve vos compétences en matière de vocabulaire et de devinettes tout en vous amusant.

Pour exécuter le jeu en local, suivez ces étapes :

1/Assurez-vous d'avoir VS installé sur votre machine.
2/Clonez le dépôt git "HangmanWeb" sur votre système local.
3/Ensuite, exécutez la commande "go run ." pour lancer le serveur :


4/Ouvrez votre navigateur et accédez à l'URL suivante : http://localhost:8080/templates/menu


Les règles du jeu sont simples :

Vous pouvez proposer une lettre ou un mot.
Si la lettre proposée ou le mot sont présents dans le mot, ces derniers seront révélés.
Si la lettre proposée n'est pas présente dans le mot, le dessin du pendu avancera d'une étape.
Si le mot proposé n'est pas le mot cherché, le dessin du pendu avancera de deux étapes.
Si vous parvenez à deviner le mot avant que le pendu ne soit complet, vous gagnez la partie. Sinon, vous perdez.

