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