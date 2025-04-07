# 🐦 Chirpy

Chirpy est une API JSON développée en Go et inspirée de plateformes comme Bluesky ou Mastodon. Ce projet utilise des outils minimalistes et intègre des fonctionnalités telles que les middlewares, le routage, la journalisation, les webhooks, l'authentification, l'autorisation, les JWTs et une base de données Postgres.

Ce projet a été réalisé dans le cadre du cours de Boot.dev sur la création d'un serveur HTTP en Go.

## 📚 Sommaire

- [⚙️ Installation](#️-installation)
- [🤖 API](#-api)

## ⚙️ Installation

### 🛠️ Étapes d'installation

1. Clonez le dépôt :

   ```bash
   git clone https://github.com/benKapl/chirpy.git
   cd chirpy
   ```

2. Construisez le projet :

   ```bash
   go build -o chirpy
   ```

3. Lancez le serveur :
   ```bash
   ./chirpy
   ```

### 🔧 Configuration

#### 🗄️ Initialisation de la base de données Postgres

1. Installez Postgres sur votre machine.
2. Créez une base de données nommée `chirpy` :
   ```sql
   CREATE DATABASE chirpy;
   ```
3. Appliquez les schémas SQL situés dans le dossier `sql/schema/` pour configurer les tables nécessaires :
   ```bash
   psql -U <votre_utilisateur> -d chirpy -f sql/schema/001_users.sql
   psql -U <votre_utilisateur> -d chirpy -f sql/schema/002_chirps.sql
   psql -U <votre_utilisateur> -d chirpy -f sql/schema/003_users.sql
   psql -U <votre_utilisateur> -d chirpy -f sql/schema/004_refresh_tokens.sql
   psql -U <votre_utilisateur> -d chirpy -f sql/schema/005_users.sql
   ```

#### 🔑 Configuration des secrets dans `.env`

Créez un fichier `.env` à la racine du projet avec les variables suivantes :

```env
DB_URL=postgres://<utilisateur>:<mot_de_passe>@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=<chaîne_aléatoire>
POLKA_KEY=f271c81ff7084ee5b99a5091b42d486e
```

- **DB_URL** : Chaîne de connexion à votre base de données Postgres locale.
- **PLATFORM** : Définissez sur `dev`.
- **JWT_SECRET** : Une chaîne aléatoire utilisée pour valider les tokens JWT.
- **POLKA_KEY** : Clé pour simuler un fournisseur de paiement. Utilisez `f271c81ff7084ee5b99a5091b42d486e` pour que cela fonctionne.\*

### 🔨 Tests

```bash
go test ./...
```

## 🤖 API

### 📍 Endpoints

Voici la liste des endpoints disponibles et leurs fonctionnalités :

#### 👤 Utilisateurs

- **POST /api/users** : Crée un nouvel utilisateur.
- **PUT /api/users** : Met à jour les informations d'un utilisateur (authentification requise).
- **POST /api/login** : Authentifie un utilisateur et retourne un token JWT.
- **POST /api/refresh** : Génère un nouveau token JWT à partir d'un token de rafraîchissement.
- **POST /api/revoke** : Révoque un token de rafraîchissement.

#### 🐤 Chirps

- **GET /api/chirps** : Récupère tous les chirps ou ceux d'un utilisateur spécifique.
  - **Paramètres de requête** :
    - `userId` (optionnel) : Filtre les chirps par l'ID de l'utilisateur.
    - `sort` (optionnel) : Définit l'ordre de tri (`asc` ou `desc`).
- **GET /api/chirps/{id}** : Récupère un chirp spécifique par son ID.
- **POST /api/chirps** : Crée un nouveau chirp (authentification requise).
- **DELETE /api/chirps/{id}** : Supprime un chirp (authentification requise).

#### 🔒 Administration

- **GET /admin/metrics** : Affiche les métriques d'utilisation.
- **POST /admin/reset** : Réinitialise les données (uniquement en environnement `dev`).

#### 🔔 Webhooks

- **POST /api/polka/webhooks** : Gère les webhooks pour les événements de mise à niveau utilisateur.

#### ❤️ Santé

- **GET /api/healthz** : Vérifie l'état de santé du serveur.

### 🛡️ Middlewares

Voici la liste des middlewares utilisés dans le projet :

1. **middlewareLog** : Journalise chaque requête HTTP avec sa méthode et son chemin.
2. **middlewareMetricsInc** : Incrémente un compteur pour suivre les visites des fichiers statiques.
3. **AuthenticateAccessToken** : Vérifie et valide les tokens JWT pour les endpoints nécessitant une authentification.

Chirpy est conçu pour être extensible et facile à utiliser, tout en respectant les bonnes pratiques de développement d'API modernes.
