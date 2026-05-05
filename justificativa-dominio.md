# Justificativa da escolha do domínio

Escolhi desenvolver um sistema de gestão de clubes do livro porque é um tema simples de entender, mas que permite representar um domínio com relacionamentos e regras de negócio reais.

Nesse contexto, uma pessoa pode criar um clube, participar de clubes existentes, sugerir livros, sortear o tema do mês e marcar encontros. Isso permite trabalhar com entidades como usuários, clubes, membros, sugestões, temas mensais e reuniões, mostrando bem a relação entre os dados no banco.

O domínio também permite aplicar diferentes níveis de permissão de forma natural. Um administrador pode gerenciar o clube, membros podem participar e sugerir livros, e visitantes podem ter acesso mais limitado. Assim, a aplicação consegue demonstrar autenticação, autorização por roles e proteção de rotas.

Além disso, o projeto permite usar DTOs para separar os dados internos da aplicação dos dados retornados pela API, evitando expor informações sensíveis e mantendo a estrutura mais organizada.

Por esses motivos, o domínio de clubes do livro foi escolhido por ser acessível, coerente com os requisitos da atividade e adequado para demonstrar persistência, relacionamentos, segurança e arquitetura em camadas.