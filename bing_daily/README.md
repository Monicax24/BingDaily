# Bing Daily

Bing Daily is a BeReal-like social media app designed for Binghamton University students, encouraging daily photo posts in response to timed community prompts. It uses Google Sign-In restricted to `@binghamton.edu` accounts and is built with Flutter for a cross-platform experience on iOS, Android, and web.

## Features
- **Google Sign-In**: Authenticate with Firebase using Google OAuth, restricted to `@binghamton.edu` accounts.
- **Community Selection**: Join preset communities (e.g., "BingLife") to share daily moments.
- **Daily Posts**: Capture/upload photos in response to timed prompts, viewable by others in the community.
- **Liking Posts**: Like posts within your selected community.
- **Responsive UI**: Consistent experience across mobile and web, following Binghamton branding (green/white theme).

## Tech Stack
- **Frontend**: Flutter (Dart) for cross-platform development.
- **State Management**: Riverpod for reactive state handling.
- **Navigation**: Go Router for declarative routing.
- **Authentication**: Firebase Authentication with Google Sign-In (v7).
- **HTTP Requests**: Dio for backend API integration.
- **Local Storage**: get_storage for caching user preferences and community data (planned).
- **Push Notifications**: Firebase Cloud Messaging (planned).

## Setup Instructions
1. **Prerequisites**:
   - Flutter SDK (v3.9.2 or higher)
   - Dart SDK
   - Firebase CLI (for Firebase configuration)
   - Android Studio/Xcode for mobile development
   - Google Cloud project with OAuth credentials

2. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-repo/bing-daily.git
   cd bing-daily
   ```

3. **Install Dependencies**:
   ```bash
   flutter pub get
   ```

4. **Configure Firebase**:
   - Set up a Firebase project and enable Google Sign-In.
   - Add `google-services.json` (Android) and `GoogleService-Info.plist` (iOS) to the respective platform directories.
   - Update `lib/firebase_options.dart` with your Firebase configuration.

5. **Run the App**:
   ```bash
   flutter run
   ```

## Project Structure
- `lib/api/`: Authentication and backend API logic (e.g., `auth_api.dart`).
- `lib/constants.dart`: Centralized colors and text constants.
- `lib/providers/`: Riverpod providers for state management (e.g., `auth_provider.dart`).
- `lib/screens/`: UI screens (e.g., `login_screen.dart`, `home_page.dart`).
- `lib/utils/`: Utility functions (e.g., `route_utils.dart` for Go Router setup).
- `lib/widgets/`: Reusable UI components (e.g., `main_screen.dart`).
- `pubspec.yaml`: Project dependencies and configuration.

## Development Guidelines
- **File Size**: Keep files under 150 lines for readability.
- **Constants**: Store colors, text, and other constants in `lib/constants.dart`.
- **Comments**: Add comments to functions for clarity, especially for team collaboration.
- **Git Workflow**:
  - Use feature branches (e.g., `feature/login-screen`).
  - Write descriptive commit messages (e.g., `feat: add community selection UI`).
  - Submit PRs for review before merging to `main`.
- **Testing**: Write unit tests for widgets and logic, integration tests for auth flow.

## Team
- **Maya**: UI/UX design (Figma).
- **Michael**: Flutter UI implementation.
- **Eland**: Functional logic, state management, navigation, storage.
- **Ryan**: UI implementation support and testing.

## Milestones
- **MVP (Nov 8, 2025)**: Login, community selection, post creation/viewing, liking.
- **Soft Deadline (Nov 21, 2025)**: UI polish, tests, bug fixes.
- **Hard Deadline (Dec 5, 2025)**: Deploy to App Store, Play Store, and web.

## Stretch Goals
- Friending system
- User-created communities
- Post captions and text responses
- Reporting/moderation features

## Dependencies
- Backend APIs for communities, prompts, posts, and likes (provided by backend team).
- Firebase setup for authentication and notifications.
- Figma designs for UI implementation.

## Contributing
- Create a feature branch for each task.
- Follow the PRD (see `PRD for Frontend of Bing Daily App.md`) for requirements.
- Use Slack for team communication and GitHub Issues for task tracking.

## License
This project is private and intended for educational use at Binghamton University.
