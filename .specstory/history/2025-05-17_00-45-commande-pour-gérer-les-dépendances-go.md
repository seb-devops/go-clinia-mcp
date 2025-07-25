<!-- Generated by SpecStory -->

# Commande pour gérer les dépendances Go (2025-05-16 20:45:19)

_**User**_

quelle est la commande pour aller chercher mes dépendance dans mon projet 

---

_**Assistant**_

Pour aller chercher (installer) les dépendances dans ton projet Go, utilise la commande suivante à la racine de ton projet :

```sh
go mod tidy
```

Cette commande va :
- Télécharger toutes les dépendances nécessaires listées dans ton fichier go.mod.
- Ajouter les dépendances manquantes et supprimer celles qui ne sont plus utilisées.

Si tu veux simplement télécharger les dépendances (sans modifier go.mod ou go.sum), tu peux aussi utiliser :

```sh
go mod download
```

Mais dans la plupart des cas, go mod tidy est la commande recommandée pour s'assurer que tout est propre et à jour.

**Résumé :**
- Place-toi à la racine de ton projet (là où se trouve go.mod)
- Exécute :  
  ```sh
  go mod tidy
  ```

Dis-moi si tu veux la commande exacte à exécuter dans ton terminal !

---

