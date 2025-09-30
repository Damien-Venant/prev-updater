[![Build Prev-Updater](https://github.com/Damien-Venant/prev-updater/actions/workflows/build-prev-updater.yml/badge.svg)](https://github.com/Damien-Venant/prev-updater/actions/workflows/build-prev-updater.yml)
[![codecov](https://codecov.io/github/Damien-Venant/prev-updater/graph/badge.svg?token=GJGMIPCGL7)](https://codecov.io/github/Damien-Venant/prev-updater)


# 🛠 prev-updater

**prev-updater** est un outil en **CLI** écrit en **Go** permettant de mettre à jour la **version prévisionnelle** d’une *User Story* sur **Azure DevOps (ADO)**.  
Il facilite l’automatisation et la standardisation de la gestion des tickets.

---

## 🚀 Installation

### 📦 Depuis les sources
```bash
git clone https://github.com/<ton-user>/prev-updater.git
cd prev-updater
go build -o prev-updater
```


### 📦 Depuis go install
```bash
go install github.com/Damien-Venant/prev-updater@latest
```

### 📦 Depuis une release
**👉 Remplacez ``VERSION`` par la version souhaitée.**
```bash
curl -o prev-updater -L https://github.com/Damien-Venant/prev-updater/releases/download/VERSION/prev-updater-linux-VERSION
chmod +x prev-updater
```

---

## ⚡ Utilisation

Utilisation de cette commande pour plus d'information

### Afficher l'aide

````bash
prev-updater help
````

### Version de l'outil

Récupération de la version de l'outil
````bash
prev-updater version
````

### Mettre a jour la prévisionelle d'un ticket

Changement de la prévisionelles des tickets :
````bash
prev-updater start -b "YOUR_BASE_URL" \
    -o "YOUR_ORGANISATION" \
    -t "YOUR_ADO_TOKEN" \
    -p "YOUR_ADO_PROJECT"
````

---

## 📜 Logs
Par défaut, les logs sont enregistrés dans :
````bash
$HOME/pre-updater/prev-updater.log
````
Ce fichier contient l’historique d’exécution et les éventuelles erreurs.
Il peut être utile lors du débogage ou pour auditer l’utilisation de l’outil.


## 🤝 Contribution

Les contributions sont les bienvenues ! Que ce soit pour corriger un bug, ajouter une fonctionnalité ou améliorer la documentation.

### 🔧 Comment contribuer

1. **Cloner** le projet
2. **Installer** l'outil lefthook puis lancer la commande suivante :
````bash
lefthook install
````
3. Créer une issue & une nouvelle branche depuis la branche `dev`.
4. Développer votre fonctionnalité/correctif.
5. Commiter vos changements en respectant [convention de commit](https://www.conventionalcommits.org/en/v1.0.0/).
6. Poussez votre branche.
7. Créer un pull request.

## 📌 Règles de contribution
 * Respecter le style de code Go (formatage avec go fmt)
 * Documenter les nouvelles fonctionnalités dans le README
 * Ajouter des tests unitaires
