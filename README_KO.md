# GO-POSTOFFICE

[中文版](README_CN.md) | [English](README.md) | [Uyghur](README_UG.md)

GO-POSTOFFICE는 Go 기반의 고성능 WebSocket 서버 구현체로, 우체국 개념 모델을 기반으로 설계되었습니다. 이 프로젝트는 주로 연결 관리, 보안 인증, 메시지(우편) 배포 및 전달을 처리합니다.

![Go-Postoffice 통신 구조](docs/imgs/global_architecture_diagram_en.png)

## 목차

1. [특징](#특징)
2. [설치](#설치)
3. [빠른 시작](#빠른-시작)
4. [구성](#구성)
5. [API 문서](#api-문서)
6. [메시지 프로토콜](#메시지-프로토콜)
7. [클라이언트 예제](#클라이언트-예제)
8. [기여하기](#기여하기)
9. [라이선스](#라이선스)

## 특징

1. **고성능 동시 처리**: Go의 goroutines와 channels를 활용하여 효율적인 동시 연결 관리를 구현합니다.
2. **유연한 메시지 라우팅**: 우체국 모델을 기반으로 점대점 및 브로드캐스트 메시지 배포를 지원하여 효율적인 메시지 전달을 실현합니다.
3. **보안 인증**: 토큰 인증을 통합하여 연결 보안을 보장합니다.
4. **구성 가능한 메시지 검증**: 메시지 형식의 정확성을 보장하기 위한 선택적 JSON Schema 검증을 지원합니다.
5. **환경 적응성**: 다양한 시나리오에서 쉽게 배포할 수 있도록 다중 환경 구성을 지원합니다.
6. **우아한 서비스 관리**: 서비스 안정성을 보장하기 위한 우아한 시작 및 종료 메커니즘을 구현합니다.
7. **확장성**: 기능 확장과 사용자 정의가 용이한 모듈식 설계를 채택했습니다.
8. **실시간 통신**: WebSocket 기반의 전이중 통신으로 실시간 데이터 교환을 지원합니다.

### 우체국 모델의 장점

- **디커플링**: 송신자와 수신자의 완전한 분리로 시스템 유연성을 향상시킵니다.
- **신뢰성**: 메시지 지속성과 재시도 메커니즘으로 안정적인 메시지 전달을 보장합니다.
- **확장성**: 새로운 메시지 유형과 처리 로직을 쉽게 추가할 수 있습니다.
- **부하 분산**: 다중 "우체국" 인스턴스를 구현하여 시스템 처리량을 증가시킬 수 있습니다.

## 기업용 AI 챗봇 통합 가이드
기업용 AI 챗봇을 빠르게 통합하는 방법에 대한 자세한 안내는 [기업용 AI 챗봇 통합 가이드](docs/enterprise_ai_chatbot_integration_guide.md)를 참조하세요.

## 설치

시스템에 Go(버전 1.23.1 이상)가 설치되어 있는지 확인하세요.

1. 저장소 복제:
   ```
   git clone https://github.com/zlz3907/go-postoffice.git
   ```

2. 프로젝트 디렉토리로 이동:
   ```
   cd go-postoffice
   ```

3. 의존성 설치:
   ```
   go mod tidy
   ```

## 빠른 시작

1. 환경 구성:
   `.env/config-dev.json`을 `.env/config-zhycit.json`으로 복사하고 필요에 따라 구성을 수정하세요.

2. 서버 실행:
   ```
   go run main.go
   ```

3. 실행 파일 빌드:

   Linux용:
   ```
   env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.env=zhycit" -o dist/go-postoffice-linux
   ```

   macOS용:
   ```
   env GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.env=zhycit" -o dist/go-postoffice-macos
   ```

   Windows용:
   ```
   env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.env=zhycit" -o dist/go-postoffice-windows.exe
   ```

   참고: `zhycit`를 원하는 환경 이름으로 변경하세요.

4. 빌드된 실행 파일 실행:

   Linux/macOS용:
   ```
   ./dist/go-postoffice-linux   # 또는 go-postoffice-macos
   ```

   Windows용:
   ```
   .\dist\go-postoffice-windows.exe
   ```

## 구성

주요 구성 항목:

- `socketPort`: WebSocket 서버 포트
- `maxWebSocketConnections`: 최대 연결 수
- `dataSource`: 데이터 소스 구성(예: Redis)

자세한 구성 지침은 [구성 문서](docs/configuration.md)를 참조하세요.

## API 문서

API 사용 지침은 [API 문서](docs/api.md)를 참조하세요.

## 메시지 프로토콜

메시지 형식 및 필드 설명은 [메시지 프로토콜 문서](docs/message-protocol.md)를 참조하세요.

## 클라이언트 예제

- [Go 클라이언트 예제](examples/go-client.go)
- [Java 클라이언트 예제](examples/JavaClient.java)
- [JavaScript 클라이언트 예제](examples/js-client.js)
- [Python 클라이언트 예제](examples/python-client.py)

## 기여하기

모든 형태의 기여를 환영합니다. 프로젝트 개발에 참여하는 방법은 [기여 가이드라인](CONTRIBUTING.md)을 참조하세요.

## 라이선스

이 프로젝트는 Apache 2.0 라이선스 하에 있습니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요. 